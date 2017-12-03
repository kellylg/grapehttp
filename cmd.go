package main

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type cmd struct {
	Name  string   `json:"name"`
	Args  []string `json:"args"`
	Paths []string `json:"paths"`
}

var cmds = []string{"ls", "mkdir", "rm", "mv", "cp"}

func (s *HTTPStaticServer) hCmd(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	log.Println(string(data))
	c := cmd{}
	if err := json.Unmarshal(data, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check command
	if !canExec(c) {
		http.Error(w, fmt.Sprintf("Forbidden, command: %s not allowed", c.Name), http.StatusForbidden)
		return
	}

	//这里更安全期间，需要校验s.Root的路径，防止出现rm -rf /这种情况，但是在grapehttp启动时，已经
	//校验过，所以这里就不用再校验
	// check permission
	for _, path := range c.Paths {
		auth := s.readAccessConf(path, r)
		// check access first
		if auth.noAccess(r) {
			log.Printf("%s have no access permission", c.Name)
			http.Error(w, "access forbidden", http.StatusForbidden)
			return
		}

		if !auth.canDelete(r) {
			if c.Name == "rm" || c.Name == "mv" || c.Name == "cp" {
				log.Printf("%s have no delete permission", c.Name)
				http.Error(w, fmt.Sprintf("%s forbidden", c.Name), http.StatusForbidden)
				return
			}
		}

		if !auth.canUpload(r) {
			if c.Name == "mkdir" || c.Name == "mv" || c.Name == "cp" {
				log.Printf("%s have no delete permission", c.Name)
				http.Error(w, fmt.Sprintf("%s forbidden", c.Name), http.StatusForbidden)
				return
			}
		}
	}

	// complete Args
	// absPath完全可用的path
	// 这里加filepath.Clean是为了防止错误中出现//dir//file这样的不美观路径
	absPaths := []string{}
	for _, p := range c.Paths {
		absPaths = append(absPaths, filepath.Clean(filepath.Join(s.Root, p)))
	}

	c.Args = append(c.Args, absPaths...)
	cmd := exec.Command(c.Name, c.Args...)
	bytes, _ := cmd.CombinedOutput()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(strings.Replace(string(bytes), filepath.Clean(s.Root), "", 1)))
}

//生成权限控制文件
//默认所有人都有访问，上传，删除权限
func genGhs(dir string) error {
	j := `{"upload":true,"delete":true,"noaccess":false}`
	y, err := yaml.JSONToYAML([]byte(j))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(dir, ".ghs.yml"), y, 0644)
	if err != nil {
		return err
	}

	return nil
}

func canExec(c cmd) bool {
	for _, name := range cmds {
		if name == c.Name {
			return true
		}
	}

	return false
}

func getUsage(root string) (string, error) {
	command := exec.Command("df", "-hT", root)
	out, err := command.CombinedOutput()
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(".*%.*/.*")
	for _, line := range strings.Split(string(out), "\n") {
		if reg.MatchString(line) {
			field := strings.Fields(line)
			fstype := field[len(field)-6]
			size := field[len(field)-5]
			used := field[len(field)-4]
			avail := field[len(field)-3]
			use := field[len(field)-2]
			mount := field[len(field)-1]
			return fmt.Sprintf("Type(%s), Size(%s), Used(%s), Avail(%s), Use%%(%s), Mounted on(%s)", fstype, size, used, avail, use, mount), nil
		}
	}

	return "", nil
}
