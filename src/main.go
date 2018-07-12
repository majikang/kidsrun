package main

import (
	"flag"
	"os"
	"runtime"
	"strings"
)

var (
	c        *Config
	currPath string
	exit     chan bool
	output   string
	buildPkg string
	started  chan bool
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
	if buildPkg != "" {
		files = strings.Split(buildPkg, ",")
	}
	NewWatcher(paths, files)
	autoBuild(files)
	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}
