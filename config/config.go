package config

import (
	"bytes"
	"encoding/json"
	"github.com/alecthomas/kingpin"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
	"runtime"
	"text/template"

	"grapehttp/pkg/vinfo"
)

type Rbac struct {
	Host string `yaml:"db_host"`
	Port string `yaml:"db_port"`
	User string `yaml:"db_user"`
	Pass string `yaml:"db_pass"`
	Name string `yaml:"db_name"`
}

type Configure struct {
	Conf            *os.File `yaml:"-"`
	Addr            string   `yaml:"addr"`
	DbInit          bool     `yaml:"-"`
	DbInitForce     bool     `yaml:"-"`
	Count           bool     `yaml:"count"`
	AdminUsername   string   `yaml:"admin_username"`
	AdminPassword   string   `yaml:"admin_password"`
	AdminEmail      string   `yaml:"admin_email"`
	Root            string   `yaml:"root"`
	Rbac            Rbac     `yaml:"rbac"`
	HTTPAuth        string   `yaml:"httpauth"`
	SimpleAuth      bool     `yaml:"simpleauth"`
	Cert            string   `yaml:"cert"`
	Key             string   `yaml:"key"`
	Cors            bool     `yaml:"cors"`
	Theme           string   `yaml:"theme"`
	XHeaders        bool     `yaml:"xheaders"`
	Upload          bool     `yaml:"upload"`
	Delete          bool     `yaml:"delete"`
	NoAccess        bool     `yaml:"noaccess"`
	PlistProxy      string   `yaml:"plistproxy"`
	Title           string   `yaml:"title"`
	Debug           bool     `yaml:"debug"`
	GoogleTrackerId string   `yaml:"google-tracker-id"`
	Auth            struct {
		Type   string `yaml:"type"`
		OpenID string `yaml:"openid"`
		HTTP   string `yaml:"http"`
	} `yaml:"auth"`
}

var (
	Gcfg              = Configure{}
	defaultPlistProxy = "https://plistproxy.herokuapp.com/plist"
	defaultOpenID     = "https://some-hostname.com/openid/"

	VERSION   = "unknown"
	BUILDTIME = "unknown time"
	GITCOMMIT = "unknown git commit"
	SITE      = "https://github.com/lexkong/grapehttp"
)

func versionMessage() string {
	t := template.Must(template.New("version").Parse(`GoHTTPServer
	Version:        {{.Version}}
	Go version:     {{.GoVersion}}
	OS/Arch:        {{.OSArch}}
	Git commit:     {{.GitCommit}}
	Built:          {{.Built}}
	Site:           {{.Site}}`))
	buf := bytes.NewBuffer(nil)
	t.Execute(buf, map[string]interface{}{
		"Version":   VERSION,
		"GoVersion": runtime.Version(),
		"OSArch":    runtime.GOOS + "/" + runtime.GOARCH,
		"GitCommit": GITCOMMIT,
		"Built":     BUILDTIME,
		"Site":      SITE,
	})
	return buf.String()
}

func getVersion() string {
	v := vinfo.Get()
	marshalled, _ := json.MarshalIndent(&v, "", "  ")
	return string(marshalled)
}

func LoadConfig() (Configure, error) {
	// initial default conf
	Gcfg.Root = "./data"
	Gcfg.Addr = ":8000"
	Gcfg.Theme = "black"
	Gcfg.PlistProxy = defaultPlistProxy
	Gcfg.Auth.OpenID = defaultOpenID
	Gcfg.GoogleTrackerId = "UA-81205425-2"
	Gcfg.Title = "Go HTTP File Server"

	kingpin.HelpFlag.Short('h')
	kingpin.Version(getVersion())
	kingpin.Flag("conf", "config file path, yaml format").Short('c').FileVar(&Gcfg.Conf)
	kingpin.Flag("root", "root directory, default ./").Short('r').StringVar(&Gcfg.Root)
	kingpin.Flag("addr", "listen address, default :8000").Short('a').StringVar(&Gcfg.Addr)
	kingpin.Flag("cert", "tls cert.pem path").StringVar(&Gcfg.Cert)
	kingpin.Flag("key", "tls key.pem path").StringVar(&Gcfg.Key)
	kingpin.Flag("simpleauth", "Simple http auth or not").BoolVar(&Gcfg.SimpleAuth)
	kingpin.Flag("auth-type", "Auth type <http|openid>").StringVar(&Gcfg.Auth.Type)
	kingpin.Flag("auth-http", "HTTP basic auth (ex: user:pass)").StringVar(&Gcfg.Auth.HTTP)
	kingpin.Flag("auth-openid", "OpenID auth identity url").StringVar(&Gcfg.Auth.OpenID)
	kingpin.Flag("theme", "web theme, one of <black|green>").StringVar(&Gcfg.Theme)
	kingpin.Flag("upload", "enable upload support").BoolVar(&Gcfg.Upload)
	kingpin.Flag("delete", "enable delete support").BoolVar(&Gcfg.Delete)
	kingpin.Flag("xheaders", "used when behide nginx").BoolVar(&Gcfg.XHeaders)
	kingpin.Flag("cors", "enable cross-site HTTP request").BoolVar(&Gcfg.Cors)
	kingpin.Flag("debug", "enable debug mode").BoolVar(&Gcfg.Debug)
	kingpin.Flag("plistproxy", "plist proxy when server is not https").Short('p').StringVar(&Gcfg.PlistProxy)
	kingpin.Flag("title", "server title").StringVar(&Gcfg.Title)
	kingpin.Flag("google-tracker-id", "set to empty to disable it").StringVar(&Gcfg.GoogleTrackerId)
	kingpin.Flag("db", "init db").Short('d').BoolVar(&Gcfg.DbInit)
	kingpin.Flag("force", "force init db first drop db then rebuild it").Short('f').BoolVar(&Gcfg.DbInitForce)

	kingpin.Parse() // first parse conf

	if Gcfg.Conf != nil {
		defer func() {
			kingpin.Parse() // command line priority high than conf
		}()
		ymlData, err := ioutil.ReadAll(Gcfg.Conf)
		if err != nil {
			return Gcfg, err
		}

		err = yaml.Unmarshal(ymlData, &Gcfg)
		return Gcfg, err
	}
	return Gcfg, nil
}
