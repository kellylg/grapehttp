package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"grapehttp/config"
	. "grapehttp/lib"
	"grapehttp/models/admin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type userinfo struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

const basicScheme string = "Basic "

func (s *HTTPStaticServer) hUserAdd(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "only `admin` user have operation authority", http.StatusInternalServerError)
		return
	}

	data, _ := ioutil.ReadAll(r.Body)
	log.Println(string(data))
	user := &admin.User{}
	if err := json.Unmarshal(data, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Createtime = time.Now()
	user.Password = Pwdhash(user.Password)
	err := user.Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Success\n"))
}

func (s *HTTPStaticServer) hUserModify(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "only `admin` user have operation authority", http.StatusInternalServerError)
		return
	}

	data, _ := ioutil.ReadAll(r.Body)
	log.Println(string(data))
	user := &admin.User{}
	if err := json.Unmarshal(data, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo := admin.GetUserOnlyByUsername(user.Username)
	if user.Password != "" {
		userInfo.Password = Pwdhash(user.Password)
	}

	if user.Email != "" {
		userInfo.Email = user.Email
	}

	if user.Remark != "" {
		userInfo.Remark = user.Remark
	}

	if user.Nickname != "" {
		userInfo.Nickname = user.Nickname
	}

	if err := userInfo.Update(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Success\n"))
}

func (s *HTTPStaticServer) hUserDisable(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "only `admin` user have operation authority", http.StatusInternalServerError)
		return
	}

	data, _ := ioutil.ReadAll(r.Body)
	log.Println(string(data))
	user := struct {
		Usernames []string `json:"usernames"`
	}{}
	if err := json.Unmarshal(data, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, username := range user.Usernames {
		userInfo := admin.GetUserOnlyByUsername(username)
		userInfo.Status = 0
		if err := userInfo.Update(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("Success\n"))
}

func (s *HTTPStaticServer) hUserEnable(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "only `admin` user have operation authority", http.StatusInternalServerError)
		return
	}

	data, _ := ioutil.ReadAll(r.Body)
	log.Println(string(data))
	user := struct {
		Usernames []string `json:"usernames"`
	}{}
	if err := json.Unmarshal(data, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, username := range user.Usernames {
		userInfo := admin.GetUserOnlyByUsername(username)
		userInfo.Status = 1
		if err := userInfo.Update(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("Success\n"))
}

func (s *HTTPStaticServer) hUserDel(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "only `admin` user have operation authority", http.StatusInternalServerError)
		return
	}

	data, _ := ioutil.ReadAll(r.Body)
	log.Println(string(data))
	user := struct {
		Usernames []string `json:"usernames"`
	}{}
	if err := json.Unmarshal(data, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, username := range user.Usernames {
		userInfo, _ := admin.GetUserByUsername(username)
		if err := userInfo.Delete(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("Success\n"))
}

func (s *HTTPStaticServer) hUserGet(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "only `admin` user have operation authority", http.StatusInternalServerError)
		return
	}

	data, _ := ioutil.ReadAll(r.Body)
	log.Println(string(data))
	user := struct {
		Username string `json:"username"`
	}{}
	if err := json.Unmarshal(data, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := admin.GetUserByUsername(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*
		if err := userInfo.Delete(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/

	resp, _ := json.MarshalIndent(userInfo, "", "    ")

	w.Write(resp)
}

func (s *HTTPStaticServer) hUserSearch(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "only `admin` user have operation authority", http.StatusInternalServerError)
		return
	}

	data, _ := ioutil.ReadAll(r.Body)
	log.Println(string(data))
	user := struct {
		Username string `json:"username"`
	}{}
	if err := json.Unmarshal(data, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users := admin.SearchUserByUsername(user.Username)
	resp, _ := json.Marshal(users)
	w.Write(resp)
}

func (s *HTTPStaticServer) hUserList(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "only `admin` user have operation authority", http.StatusInternalServerError)
		return
	}

	data, _ := ioutil.ReadAll(r.Body)
	log.Println(string(data))
	num := struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}{}
	if err := json.Unmarshal(data, &num); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users := admin.ListUserByUsername(num.Offset, num.Limit)
	resp, _ := json.Marshal(users)
	w.Write(resp)
}

func grapeAuthFunc(user, pass string, r *http.Request) bool {
	userInfo, err := admin.GetUsableUserByUsername(user)
	if err != nil {
		return false
	}

	if Strtomd5(pass) != userInfo.Password {
		log.Printf("user: %s login failed", user)
		return false
	}

	if !config.Gcfg.Debug {
		return true
	}

	log.Printf("user: %s login success", user)
	// 更新登陆时间
	userInfo = admin.UpdateLoginTime(&userInfo)
	userInfo.Logincount += 1
	userInfo.Update()
	return true
}

func isAdmin(r *http.Request) bool {

	if getUser(r) == "admin" {
		return true
	}

	return false
}

func getUser(r *http.Request) string {
	// Confirm the request is sending Basic Authentication credentials.
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, basicScheme) {
		return ""
	}

	str, err := base64.StdEncoding.DecodeString(auth[len(basicScheme):])
	if err != nil {
		return ""
	}

	creds := bytes.SplitN(str, []byte(":"), 2)
	if len(creds) != 2 {
		return ""
	}

	givenUser := string(creds[0])
	//givenPass := string(creds[1])
	return givenUser
}
