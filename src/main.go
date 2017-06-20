package main

import (
	"flag"
	"fmt"
	"utils"
)

var (
	cfg      *utils.Config
	currpath string
	exit     chan bool
	output   string
	buildPkg string

	started chan bool
)

func init() {
	//定义命令行参数
	flag.StringVar(&output, "o", "", "go build output")
}

func main() {
	//加载配置文件
	var c = utils.Config{}
	c.LoadConfig()
	fmt.Printf("%v", c)
	output = "123"
	flag.Parse()
}
