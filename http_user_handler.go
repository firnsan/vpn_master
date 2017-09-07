package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"strings"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
	"time"
	"strconv"
)

func init() {
	gHttpServer.AddToInit(InitUserHandler)
	gHttpServer.AddToUninit(UninitUserHandler)
}

func InitUserHandler() error {
	return nil
}

func UninitUserHandler() {
}



func UserRandHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	// 管理密钥
	secret := strings.TrimSpace(r.FormValue("secret"))
	zone := strings.TrimSpace(r.FormValue("zone"))
	months := strings.TrimSpace(r.FormValue("months"))

	if secret == "" || secret != gApp.Cnf.AdminSecret  {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if zone == "" {
		http.Error(w, "empty zone", http.StatusBadRequest)
		return
	}
	if months == "" {
		months = "1"
	}
	monthsInt,err := strconv.Atoi(months)
	if err != nil {
		http.Error(w, "invalid months", http.StatusBadRequest)
		return
	}
	
	user := NewRandUser()
	user.Zone = zone
	user.CreateTime = time.Now()
	user.ExpireTime = user.CreateTime.AddDate(0,monthsInt,0)

	slaves := []Slave{{Address:"119.28.21.143:3390"}}
	for _,slave := range slaves {
		err = deployUser(slave, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}
	
	dborm := orm.NewOrm()
	_, err = dborm.Insert(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	
	buf, _ := json.Marshal(user)
	w.Write(buf)

}


func deployUser(slave Slave, user *User) error {
	method := "http://" + slave.Address + "/user/deploy"
	buf, _ := json.Marshal(user)

	log.Printf("asking %s, with request: %s\n", method, string(buf[:]))
	r := strings.NewReader(string(buf[:]))
	resp,err := http.Post(method, "text/json", r)
	if err != nil {
		log.Printf("deployUser failed: %s", err)
		return err
	}
	resp.Body.Close()

	return nil
}

