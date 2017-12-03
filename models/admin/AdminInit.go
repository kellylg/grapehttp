package admin

import (
	"github.com/astaxie/beego/orm"
	"grapehttp/config"
	"grapehttp/lib"
	"log"
	"time"
)

func InitData() {
	InsertUser()
}

// 用户数据
func InsertUser() {
	log.Println("insert user ...")
	u := new(User)
	u.Username = config.Gcfg.AdminUsername
	u.Nickname = "HttpAdmin"
	u.Password = lib.Pwdhash(config.Gcfg.AdminPassword)
	u.Email = config.Gcfg.AdminEmail
	u.Remark = "God in Grapehttp Country"
	// 2 stand for close, but it has very high authority
	u.Status = 1
	//u.Createtime = lib.GetTime()
	u.Createtime = time.Now()
	err := u.Insert()
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("insert user end")
}

func (m *User) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}
