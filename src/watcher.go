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

func NewWatcher(paths []string, files []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf(" Fail to create new Watcher[ %s ]\n", err)
		os.Exit(2)
	}

	go func() {
		for {
			select {
			case e := <-watcher.Events:
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

						autoBuild(files)
					}()
				}
			case err := <-watcher.Errors:
				log.Printf("watcher error, err:%s", err.Error())
			}
		}
	}()

	log.Printf("Initializing watcher...\n")
	for _, path := range paths {
		log.Printf("Directory( %s )\n", path)
		err = watcher.Add(path)
		if err != nil {
			log.Printf("Fail to watch directory[ %s ]\n", err)
			os.Exit(2)
		}
	}

}

func autoBuild(files []string) {
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

	bcmd := exec.Command(cmdName, args...)
	bcmd.Env = append(os.Environ(), "GOGC=off")
	bcmd.Stdout = os.Stdout
	bcmd.Stderr = os.Stderr
	err = bcmd.Run()
	if err != nil {
		log.Printf("============== Build failed ===================\n")
		return
	}
	log.Printf("Build was successful\n")
	Restart(c.Output)
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

func Start(appname string) {
	log.Printf("Restarting %s ...\n", appname)
	if strings.Index(appname, "./") == -1 {
		appname = "./" + appname
	}

	cmd = exec.Command(appname)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = append([]string{appname}, c.Cmds...)

	go cmd.Run()
	log.Printf("%s is running...\n", appname)
	started <- true
}

func Restart(appname string) {
	Kill()
	go Start(appname)
}
