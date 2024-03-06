package handle

import (
	"github.com/flosch/pongo2/v4"
	"imgGo-Server/global"
	"net/http"
)

func SyncViewLogsView(w http.ResponseWriter, r *http.Request) {
	global.MyDebug("--call SyncViewLogsView--")
	//解析模板文件
	var t = pongo2.Must(pongo2.FromFile("./template/sync_logs.html"))

	//输出文件数据
	err := t.ExecuteWriter(pongo2.Context{
		"conf":  global.Conf,
		"query": r.FormValue("query"),
	}, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
