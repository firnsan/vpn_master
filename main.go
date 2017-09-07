package main

import (
	_ "fmt"
	"github.com/firnsan/incubator"
)

var (
	gApp = NewApplication()
)

func main() {
	// 对程序进行孵化
	incubator.Incubate(gApp)

	// 开始运行
	gApp.Run()
}
