package cmd

import (
	"regexp"
	"time"

	"github.com/asaskevich/govalidator"
)

type user struct {
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Logincount    int       `json:"loginCount"`
	Repassword    string    `json:"repassword"`
	Nickname      string    `json:"nickname"`
	Email         string    `json:"email"`
	Remark        string    `json:"remark"`
	Status        int       `json:"status"`
	Lastlogintime time.Time `json:"lastLoginTime"`
	Createtime    time.Time `json:"createTime"`
	Lastip        string    `json:"lastip"`
}

const (
	TABLE_WIDTH = 80
)

func isUsername(username string) bool {
	reg := regexp.MustCompile("^[a-zA-Z]+$")
	return reg.MatchString(username)
}

func isEmail(email string) bool {
	return govalidator.IsEmail(email)
}
