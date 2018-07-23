package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	cmd          *exec.Cmd
	state        sync.Mutex
	eventTime    = make(map[string]int64)
	scheduleTime time.Time
)

func NewWatcher(paths, files []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf(" Fail to create new Watcher[ %s ]\n", err)
		os.Exit(2)
	}

	go func() {
		for {
			select {
			case e := <-watcher.Events:
				log.Printf("接收到e事件 # %+v #\n", e)

				isbuild := true

				// Skip ignored files
				if isTmpIgnoreFile(e.Name) {
					continue
				}
				if !isWatchExt(e.Name) {
					continue
				}

				mt := getFileModTime(e.Name)
				if t := eventTime[e.Name]; mt == t {
					log.Printf("[SKIP] # %s #\n", e.String())
					isbuild = false
				}

				eventTime[e.Name] = mt

				if isbuild {
					go func() {
						// Wait 1s before autobuild util there is no file change.
						scheduleTime = time.Now().Add(1 * time.Second)
						for {
							time.Sleep(scheduleTime.Sub(time.Now()))
							if time.Now().After(scheduleTime) {
								break
							}
							return
						}
						log.Printf("重新编译 # %+v #\n", files)

						autoBuild(files, false)
					}()
				}
			case err := <-watcher.Errors:
				log.Printf("watcher error, err:%s", err.Error())
			}
		}
	}()

	log.Printf("Initializing watcher......\n")
	for _, path := range paths {
		log.Printf("Directory( %s )\n", path)
		err = watcher.Add(path)
		if err != nil {
			log.Printf("Fail to watch directory[ %s ]\n", err)
			os.Exit(2)
		}
	}

}

func autoBuild(files []string, first bool) {
	state.Lock()
	defer state.Unlock()

	log.Printf("Start building...\n")

	os.Chdir(currPath)

	cmdName := "go"

	var err error

	args := []string{"build"}
	args = append(args, "-o", c.Output)
	if c.BuildTags != "" {
		args = append(args, "-tags", c.BuildTags)
	}
	args = append(args, files...)

	cmd = exec.Command(cmdName, args...)
	cmd.Env = append(os.Environ(), "GOGC=off")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("================ Build failed ===================\n")
		return
	}
	log.Printf("Build successful!!")
	Restart(first)
}

func Start() {
	log.Printf("server %s is starting...\n", c.AppName)
	if strings.Index(c.Output, "/") == -1 {
		c.Output = "./" + c.Output
	}

	cmd = exec.Command(c.Output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = append([]string{c.Output}, c.Cmds...)

	go cmd.Run()
	log.Printf("%s running...\n", c.AppName)
}

func Restart(first bool) {
	if !first {
		Kill()
	}
	Start()
}

func Kill() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Kill.recover -> ", e)
		}
	}()

	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println("Kill -> ", err)
		}
	}
}
