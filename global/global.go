package global

import (
	"imgGo-Server/model"
	"log"
)

var Conf *model.ConfigModel

// MyDebug 自定义简易日志输出
func MyDebug(a ...any) {
	if Conf.IsDebug {
		log.Println(a...)
	}
}
