package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"path/filepath"
)

type Config struct {
	Appname       string   `yaml:"appname"`
	Output        string   `yaml:"output"` //指定ouput执行的程序路径
	WatchFileExts []string `yaml:"watch_file_exts"`
	Cmds          []string `yaml:"cmds"`
	ExcludedPaths []string `yaml:"excluded_paths"`
}

func (c *Config) LoadConfig() *Config {
	var configFile = "./watch.yml"
	//检测地址是否是绝对地址，是绝对地址直接返回，不是绝对地址，会添加当前工作路径到参数path前，然后返回
	filename, _ := filepath.Abs(configFile)
	//判断文件是否存在
	if !fileExist(filename) {
		return c
	}
	//读取yaml文件
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	//格式化配置文件到yamlConfig
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic(err)
	}
	//根据实际开发环境加载配置信息
	c.LoadLocal()
	return c
}

func (c *Config) LoadLocal() *Config {
	/*currpath, _ = os.Getwd() //获取当前路径
	//app名默认取当前文件夹项目名或者文件名，如果有output则换成output的
	if output == "" {
		cfg.AppName = path.Base(currpath)
	} else {
		cfg.AppName = path.Base(output)
	}*/

	return c
}
