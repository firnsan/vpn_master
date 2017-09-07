package main

import (
	log "github.com/Sirupsen/logrus"
	"time"
	"github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql" // import your required driver
)

var (
	VERSION = "0.1.0"
)

type Application struct {
	Cnf     *Config
	cnfPath string
}

func NewApplication() *Application {
	return &Application{cnfPath: "../etc/app.cnf"}
}

func (o *Application) GetVersion() string {
	return VERSION
}

func (o *Application) GetUsage() string {
	return `Tar Stream Server.

	Usage:
	tar_stream [--cnf=<path>]
	tar_stream -h | --help
	tar_stream --version

	Options:
	--cnf=<path>  config file path [default: ../etc/app.cnf].`

}

func (o *Application) OnOptParsed(m map[string]interface{}) {
	o.cnfPath = m["--cnf"].(string)
}

func (o *Application) OnReload() error {
	var err error
	log.Warn("application need to reload")

	// 重载配置
	if o.Cnf != nil {
		err = o.Cnf.Reload()
	}

	// 重载日志
	UninitLog()
	err = InitLog()
	if err != nil {
		return err
	}

	log.Warn("application reloaded")

	return nil
}

func (o *Application) OnStop() {
	log.Warn("application need to stop")
	gHttpServer.Stop()
	// 最多只苟延1秒, 目的让主线程退出
	time.Sleep(time.Second)
	log.Warn("application stopped")
}

func (o *Application) Run() {
	var err error

	// 加载配置
	o.Cnf, err = NewConfig(o.cnfPath)
	if err != nil {
		log.Fatalf("init config failed: %s", err)
	}

	// 初始化日志
	err = InitLog()
	if err != nil {
		log.Fatalf("init log failed: %s", err)
	}
	defer UninitLog()

	log.Print("")
	log.Print("application started")


	// 初始化orm
	orm.Debug = true
	err = orm.RegisterDataBase("default", "mysql", o.Cnf.DbConnect)
	if err != nil {
		log.Fatalf("init database connection failed: %s", err)
	}

	// 启动http server
	err = gHttpServer.Run()
	if err != nil {
		log.Fatalf("application quited, because http server quited abnormally: %s", err)
	}

	log.Warn("application quited")

}
