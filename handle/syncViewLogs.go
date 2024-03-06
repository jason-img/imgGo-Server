package handle

import (
	"fmt"
	"imgGo-Server/dao"
	"imgGo-Server/global"
	"imgGo-Server/util"
	"net/http"
)

// SyncViewLogs 处理同步文件访问记录
func SyncViewLogs(w http.ResponseWriter, r *http.Request) {
	conf := global.Conf

	global.MyDebug("--call SyncViewLogs--")

	if r.Method != "POST" {
		SyncViewLogsView(w, r)
		return
	}

	var lineCount = 0
	var recordCount = 0
	lastMd5 := util.GetLastFileMd5()

	global.MyDebug("lastMd5 ->", lastMd5)

	nginxAccessLogMd5 := util.GetFileMd5(conf.NginxAccessLogPath)

	global.MyDebug("nginxAccessLogMd5 ->", nginxAccessLogMd5)

	if lastMd5 != "" && lastMd5 == nginxAccessLogMd5 {
		util.ResponseComm(w, 200, "日志文件未更新，跳过")
		return
	}
	lineCount, recordCount = dao.ReadLogFile(conf.NginxAccessLogPath)
	lastMd5 = util.GetFileMd5(conf.NginxAccessLogPath)
	util.UpdateLastFileMd5(lastMd5)
	util.ResponseComm(w, 200, fmt.Sprintf("记录日志成功，共处理：%d行，记录：%d行", lineCount, recordCount))
}
