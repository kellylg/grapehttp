/*
	Copyright 2017 by GoWeb author: gdccmcm14@live.com.
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at
		http://www.apache.org/licenses/LICENSE-2.0
	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License
*/
package admin

import (
	"time"

	"github.com/astaxie/beego/orm"
	. "grapehttp/lib"
	//_ "grapehttp/models"
)

type User struct {
	Id            int64
	Logincount    int       `orm:"column(logincount)" form:"logincount" json:"loginCount"`
	Username      string    `orm:"unique;size(32)" form:"Username"  valid:"Required;MaxSize(20);MinSize(6)" json:"username"`
	Password      string    `orm:"size(32)" form:"Password" valid:"Required;MaxSize(20);MinSize(6)" json:"password"`
	Repassword    string    `orm:"-" form:"Repassword" valid:"Required" json:"repassword"`
	Nickname      string    `orm:"unique;size(32)" form:"Nickname" valid:"Required;MaxSize(20);MinSize(2)" json:"nickname"`
	Email         string    `orm:"size(32)" form:"Email" valid:"Email" json:"email"`
	Remark        string    `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)" json:"remark"`
	Status        int       `orm:"default(2)" form:"Status" valid:"Range(0,1)" json:"status"`
	Lastlogintime time.Time `orm:"null;type(datetime)" form:"-" json:"lastLoginTime"`
	Createtime    time.Time `orm:"type(datetime)" json:"createTime"`
	Lastip        string    `json:"lastip"`
}

func (u *User) TableName() string {
	return "tb_http_user"
}

func init() {
	orm.RegisterModel(new(User))
}

func GetUserByUsername(username string) (user User, err error) {
	o := orm.NewOrm()
	o.Using("caj")

	user = User{Username: username}
	err = o.Read(&user, "Username")
	return user, err
}

func GetUsableUserByUsername(username string) (user User, err error) {
	o := orm.NewOrm()
	o.Using("caj")

	user = User{Username: username, Status: 1}
	err = o.Read(&user, "Username", "Status")
	return user, err
}

func GetUserOnlyByUsername(username string) (user User) {
	o := orm.NewOrm()
	o.Using("caj")

	user = User{Username: username}
	o.Read(&user, "Username")
	return user
}

func SearchUserByUsername(username string) []*User {
	o := orm.NewOrm()
	o.Using("caj")

	qs := o.QueryTable("tb_http_user")
	qs = qs.Filter("username__icontains", username)
	//qs = qs.Filter("status", 1)
	users := make([]*User, 0)
	qs.All(&users)

	return users
}

func ListUserByUsername(offset, limit int) []*User {
	o := orm.NewOrm()
	o.Using("caj")

	qs := o.QueryTable("tb_http_user")
	qs = qs.Limit(limit, offset)
	//qs = qs.Filter("status", 1)
	users := make([]*User, 0)
	qs.All(&users)

	return users
}

func (m *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func UpdateLoginTime(u *User) User {
	u.Lastlogintime = GetTime()
	o := orm.NewOrm()
	o.Update(u)
	return *u
}

func (m *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}
