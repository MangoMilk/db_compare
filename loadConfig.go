package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//解析yml文件
type Config struct {
	APP_NAME string `yaml:"APP_NAME"`
	DBA      DBConf `yaml:"DBA"`
	DBB      DBConf `yaml:"DBB"`
}

type DBConf struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

var conf *Config

func (c *Config) GetConf() *Config {

	// configPath, _ := os.Getwd()
	goPath := os.Getenv("GOPATH")
	var configPath string = goPath + "/src/db_compare/config.yaml"

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic("配置读取错误: " + err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic("配置解析错误: " + err.Error())
	}
	return c
}

func InitConfig() {
	var c Config
	conf = c.GetConf()
}
