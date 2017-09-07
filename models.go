package main

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModelWithPrefix("tb_", new(User))
}
