package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var (
	c        *Config
	currPath string
	output   string
	sign     chan os.Signal
	paths    []string
)

var ignoredFilesRegExps = []string{
	`.#(\w+).go`,
	`.(\w+).go.swp`,
	`(\w+).go~`,
	`(\w+).tmp`,
}

func init() {
	//定义命令行参数
	flag.StringVar(&output, "o", "", "go build output")
}

func main() {
	//Parse cmd arguments
	flag.Parse()
	var err error
	c, err = getConfig()
	if err != nil {
		panic(err)
	}

	if c.WatchPath == "" || c.WatchPath == "./" {
		//Get the current path
		currPath, _ = os.Getwd()
	} else {
		currPath = c.WatchPath
	}

	readDirectory(currPath, &paths, c)

	runApp()

}

func runApp() {
	files := []string{}

	NewWatcher(paths, files)
	autoBuild(files, true)
	sign = make(chan os.Signal)
	signal.Notify(sign, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	<-sign
}
