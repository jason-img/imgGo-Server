package config

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"imgGo-Server/model"
	"os"
	"strings"
)

var conf *model.ConfigModel

// InitConfig 解析并初始化配置
func InitConfig() *model.ConfigModel {

	filepath := "config.yml"

	if fileSuffix := getEnv(); fileSuffix != "" {
		filepath = "config_" + fileSuffix + ".yml"
	}

	yamlBytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("YAML配置文件读取失败 #%v ", err)
		panic(err)
	}

	//fmt.Println("config file text ->", string(yamlBytes))

	err = yaml.Unmarshal(yamlBytes, &conf)
	if err != nil {
		panic(err)
	}

	//confJsonBytes, _ := json.Marshal(conf)
	//fmt.Println("Current conf", filepath, "-->", string(confJsonBytes))

	return conf
}

// getEnv 获取环境变量
func getEnv() string {
	var envBytes []byte
	var err error
	var reader *bufio.Reader
	var f *os.File
	fileSuffix := ""

	f, err = os.Open(".env")
	if err != nil {
		fmt.Println("打开.env 失败", err)
		goto RE
	}

	reader = bufio.NewReader(f)
	envBytes, _, err = reader.ReadLine()
	if err != nil {
		fmt.Println("读取.env 失败", err)
		goto RE
	}

	if strings.ToLower(string(envBytes)) != "prd" {
		fileSuffix = strings.ToLower(string(envBytes))
	}

RE:
	return fileSuffix
}
