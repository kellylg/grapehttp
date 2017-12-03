// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"grapehttp/config"
	"grapehttp/models/admin"
	"log"
	"os"
)

var (
	defaultMaxIdleConns   = 100
	defaultConnectTimeout = "10s"
	defaultReadTimeout    = "30s"
	defaultWriteTimeout   = "60s"
)

func Run() {
	if config.Gcfg.SimpleAuth && (config.Gcfg.DbInit || config.Gcfg.DbInitForce) {
		log.Fatal("db init not support with `--simpleauth` or set `simpleauth: true` in configure file")

	}
	if config.Gcfg.SimpleAuth {
		return
	}

	if config.Gcfg.DbInit {
		Syncdb(config.Gcfg.DbInitForce)
		os.Exit(0)
	}

	Connect()
}

func Connect() {
	log.Printf("database start to connect")
	rbac := config.Gcfg.Rbac
	var dns string
	db_type := "mysql"
	db_host := rbac.Host
	db_port := rbac.Port
	if db_port == "" {
		db_port = "3306"
	}
	db_user := rbac.User
	db_pass := rbac.Pass
	db_name := rbac.Name
	switch db_type {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DRMySQL)
		//orm.DefaultTimeLoc = time.UTC
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
		break
	default:
		log.Printf("db driver not support: %s", db_type)
		return
	}

	err := orm.RegisterDataBase("default", db_type, dns)
	if err != nil {
		log.Printf("register data:" + err.Error())
		panic(err.Error())
	}

	err = orm.RegisterDataBase(rbac.Name, db_type, dns)
	if err != nil {
		log.Printf("register data:" + err.Error())
		panic(err.Error())
	}
}

func Createtb() {
	admin.InitData()
}

func Syncdb(force bool) {
	//Createdb(force)
	//安全起见禁止drop database
	Createdb(false)
	Connect()
	Createconfig()
	Createtb()
	log.Printf("sync db end, please reopen app again")
}

func Createconfig() {
	name := config.Gcfg.Rbac.Name // database alias name
	force := true                 // drop table force
	verbose := true               // print log
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		log.Fatalf("database config set to force error:%s", err.Error())
	}
}

//创建数据库
func Createdb(force bool) {
	rbac := config.Gcfg.Rbac

	db_type := "mysql" // current only support mysql
	db_host := rbac.Host
	db_port := rbac.Port
	db_user := rbac.User
	db_pass := rbac.Pass
	db_name := rbac.Name

	var dns string
	var sqlstring, sql1string string

	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", db_user, db_pass, db_host, db_port)
	sqlstring = fmt.Sprintf("CREATE DATABASE if not exists `%s` CHARSET utf8 COLLATE utf8_general_ci", db_name)
	sql1string = fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", db_name)
	if force {
		fmt.Println(sql1string)
	}
	fmt.Println(sqlstring)

	db, err := sql.Open(db_type, dns)
	if err != nil {
		panic(err.Error())
	}
	if force {
		_, err = db.Exec(sql1string)
	}
	_, err1 := db.Exec(sqlstring)
	if err != nil || err1 != nil {
		log.Printf("db exec error：%v, %v", err, err1)
		panic(err.Error())
	} else {
		log.Printf("database %s created", db_name)
	}
	defer db.Close()
	log.Printf("create database end")
}
