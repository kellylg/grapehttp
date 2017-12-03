package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"grapehttp/config"
	"grapehttp/models"
	"grapehttp/pkg/vinfo"

	"github.com/go-yaml/yaml"
	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	accesslog "github.com/mash/go-accesslog"
)

type logger struct{}

func (l logger) Log(record accesslog.LogRecord) {
	log.Printf("%s - %s %d %s", record.Ip, record.Method, record.Status, record.Uri)
}

var (
	l                 = logger{}
	defaultPlistProxy = "https://plistproxy.herokuapp.com/plist"
	VERSION           = "unknown"
	blackPath         = []string{"", "/", "/bin", "/boot", "/dev", "/etc", "/home", "/lib", "/lib64", "/media", "/mnt", "/opt", "/proc", "/root", "/run", "/srv", "/sys", "/var", "/usr", "/data", "/tmp"}
)

func main() {
	gcfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if gcfg.Conf == nil && !gcfg.SimpleAuth {
		log.Fatal("no configire file specified, you need to specify `--simpleauth` to avoid database connection")
	}

	// db init
	models.Run()

	if gcfg.Debug {
		data, _ := yaml.Marshal(gcfg)
		fmt.Printf("--- config ---\n%s\n", string(data))
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// check Root first
	if !canRoot(gcfg.Root) {
		log.Fatal(fmt.Errorf("can not use root: %s", gcfg.Root))
	}
	ss := NewHTTPStaticServer(gcfg.Root)
	ss.Theme = gcfg.Theme
	ss.Title = gcfg.Title
	ss.GoogleTrackerId = gcfg.GoogleTrackerId
	ss.Upload = gcfg.Upload
	ss.Delete = gcfg.Delete
	ss.NoAccess = gcfg.NoAccess
	ss.AuthType = gcfg.Auth.Type
	usage, err := getUsage(gcfg.Root)
	if err != nil {
		log.Fatal(fmt.Errorf("can not get root usage: %v", err))
	}
	//ss.Usage = strings.Split(string(usage), "\n")[1]
	ss.Usage = string(usage)
	v := vinfo.Get()
	marshalled, _ := json.Marshal(&v)
	ss.Version = string(marshalled)

	if gcfg.PlistProxy != "" {
		u, err := url.Parse(gcfg.PlistProxy)
		if err != nil {
			log.Fatal(err)
		}
		u.Scheme = "https"
		ss.PlistProxy = u.String()
	}

	var hdlr http.Handler = ss

	hdlr = accesslog.NewLoggingHandler(hdlr, l)

	// HTTP Basic Authentication
	userpass := strings.SplitN(gcfg.Auth.HTTP, ":", 2)
	switch gcfg.Auth.Type {
	case "http":
		if len(userpass) == 2 {
			user, pass := userpass[0], userpass[1]

			if gcfg.SimpleAuth {
				hdlr = httpauth.SimpleBasicAuth(user, pass)(hdlr)
			} else {
				authOpts := httpauth.AuthOptions{
					Realm:    "Restricted",
					User:     user,
					Password: pass,
					AuthFunc: grapeAuthFunc,
					//UnauthorizedHandler: hdlr,
				}
				hdlr = httpauth.BasicAuth(authOpts)(hdlr)
			}
		}
	case "openid":
		handleOpenID(false) // FIXME(ssx): set secure default to false
	}
	// CORS
	if gcfg.Cors {
		hdlr = handlers.CORS()(hdlr)
	}
	if gcfg.XHeaders {
		hdlr = handlers.ProxyHeaders(hdlr)
	}

	http.Handle("/", hdlr)
	http.HandleFunc("/-/sysinfo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(map[string]interface{}{
			"version": VERSION,
		})
		w.Write(data)
	})

	if !strings.Contains(gcfg.Addr, ":") {
		gcfg.Addr = ":" + gcfg.Addr
	}
	log.Printf("listening on %s\n", strconv.Quote(gcfg.Addr))

	if gcfg.Key != "" && gcfg.Cert != "" {
		err = http.ListenAndServeTLS(gcfg.Addr, gcfg.Cert, gcfg.Key, nil)
	} else {
		err = http.ListenAndServe(gcfg.Addr, nil)
	}
	log.Fatal(err)
}

//测试场景：
// ., .., /, /bin, //, //bin//, ""
func canRoot(root string) bool {
	abs, _ := filepath.Abs(filepath.Clean(root))
	for _, p := range blackPath {
		if abs == p {
			return false
		}
	}

	return true
}
