package main

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AppName         string   `yaml:"appname"` //应用名称 如果是文件路径则是文件名，如果是文件夹则是文件夹名
	WatchPath       string   `yaml:"watch_path"`
	Output          string   `yaml:"output"`            //程序输出路径
	WatchFileExts   []string `yaml:"watch_file_exts"`   //监听的文件后缀名
	Cmds            []string `yaml:"cmds"`              //命令
	Excluded        []string `yaml:"excluded"`          //不需要监控的目录或者文件
	IgnoredFileExts []string `yaml:"ignored_file_exts"` //忽略的文件后缀名
	BuildTags       string   `yaml:"build_tags"`        //在go build 时期接收的-tags参数
	VendorWatch     bool     `yaml:"vendor_watch"`      //vendor 目录下的文件是否也监听
}

var configFile = "watch.yml"

func (c *Config) NewConfig() (err error) {
	//检测地址是否是绝对地址，是绝对地址直接返回，不是绝对地址，会添加当前工作路径到参数path前，然后返回
	filename, _ := filepath.Abs(configFile)
	//判断文件是否存在
	if !isExist(filename) {
		return errors.New("watch.yaml配置文件或文件夹不存在")
	}
	//读取yaml文件
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//格式化配置文件到yamlConfig
	err = yaml.Unmarshal(yamlFile, &c)

	return
}

func getConfig() (*Config, error) {
	c := &Config{}
	err := c.NewConfig()
	if err != nil {
		return c, err
	}
	if output != "" {
		c.Output = output
	}

	//app名默认取output
	if c.Output == "" {
		if c.AppName == "" {
			c.AppName = path.Base(currPath)
		}
		outputExt := ""
		if runtime.GOOS == "windows" {
			outputExt = ".exe"
		}
		c.Output = "./" + c.AppName + outputExt
	} else {
		c.AppName = path.Base(c.Output)
	}

	//监听的文件后缀
	c.WatchFileExts = append(c.WatchFileExts, ".go")

	return c, nil
}
