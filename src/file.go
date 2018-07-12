package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

// getFileModTime retuens unix timestamp of `os.File.ModTime` by given path.
func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Fail to open file[ %s ]\n", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Printf("Fail to get file information[ %s ]\n", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

// readDirectory and save filePaths to array
func readDirectory(directory string, paths *[]string, c *Config) {
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Printf("read dir failed, err:%s", err.Error())
		return
	}

	useDirectory := false
	//遍历读取目录下的文件
	for _, fileInfo := range fileInfos {
		for _, ignoredFileExt := range c.IgnoredFileExts {
			if ignoredFileExt != "" && strings.HasSuffix(fileInfo.Name(), ignoredFileExt) {
				continue
			}
		}

		if !c.VendorWatch && strings.HasSuffix(fileInfo.Name(), "vendor") {
			continue
		}

		if isExcluded(path.Join(directory, fileInfo.Name()), &c.Excluded) {
			continue
		}

		if fileInfo.IsDir() == true && fileInfo.Name()[0] != '.' {
			readDirectory(directory+"/"+fileInfo.Name(), paths, c)
			continue
		}
		if useDirectory == true {
			continue
		}
		*paths = append(*paths, directory)
		useDirectory = true
	}
	return
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// isExcluded Check if the file is excluded
func isExcluded(fileName string, excluded *[]string) bool {
	for _, p := range *excluded {
		if strings.HasSuffix(fileName, p) {
			return true
		}
	}
	return false
}

func isTmpIgnoreFile(filename string) bool {
	for _, regex := range ignoredFilesRegExps {
		r, err := regexp.Compile(regex)
		if err != nil {
			panic("Could not compile the regex: " + regex)
		}
		if r.MatchString(filename) {
			return true
		} else {
			continue
		}
	}
	return false
}

func isWatchExt(name string) bool {
	for _, s := range c.WatchFileExts {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}
