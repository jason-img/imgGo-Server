package main

import (
	"imgGo-Server/config"
	"imgGo-Server/dao"
	"imgGo-Server/global"
	"imgGo-Server/handle"
	mw "imgGo-Server/middleware"
	filter "imgGo-Server/templateFilter"
	"imgGo-Server/util"
	"log"
	"net/http"
	"os"
	"time"
)

// 删除调试信息
// go build -trimpath -ldflags="-s -w"
// 传参编译（添加版本号）
// go build -trimpath -ldflags="-s -w -X 'main.Version=v1.0.0'"
// UPX压缩
// go build -trimpath -ldflags="-s -w" -o server . && upx -9 server

var Version = "v0.6.20240131"

func init() {
	log.Println("\n\n\n\n----Staring at:", time.Now(), "----")
	log.Println("\nVersion:", Version)
	global.Conf = config.InitConfig()
	dao.InitDao()
	if global.Conf.FilerunEnable {
		dao.InitFilerunDao()
	}
}

func main() {
	conf := global.Conf

	err := util.CreateDir(conf.Path)
	if err != nil {
		panic(err)
	}
	if &conf == nil {
		os.Exit(-1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/file/upload", mw.Auth(http.HandlerFunc(handle.FileUpload)))
	mux.HandleFunc("/file/delete", mw.Auth(http.HandlerFunc(handle.FileDelete)))
	mux.HandleFunc("/syncViewLogs", http.HandlerFunc(handle.SyncViewLogs))

	if conf.FilerunEnable { // 启用Filerun 接口
		mux.HandleFunc("/filerun/", mw.Auth(http.HandlerFunc(handle.Filerun)))
	}
	if conf.UserConfig.AllowLogin { // 允许登录
		mux.HandleFunc("/user/login", mw.UserPwdCheck(http.HandlerFunc(handle.UserLogin)))
	}
	if conf.UserConfig.AllowRegister { // 允许注册
		mux.HandleFunc("/user/reg", mw.UserPwdCheck(http.HandlerFunc(handle.UserReg)))
	}

	if conf.AppMode == "dev" {
		mux.HandleFunc("/home", handle.IndexView) // TODO 完成页面登录后开放到首页
		// 加载静态文件
		mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
		mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))
		mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
		mux.Handle("/", http.FileServer(http.Dir("./upload")))
	}

	filter.Register()

	log.Println(conf.AppName+" listen on:", conf.Host+":"+conf.Port)
	log.Println("Current AppMode:", conf.AppMode)
	err = http.ListenAndServe(conf.Host+":"+conf.Port, mux)
	if err != nil {
		panic(err)
	}
}
