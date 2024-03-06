package handle

import (
	"encoding/json"
	"fmt"
	"imgGo-Server/dao"
	"imgGo-Server/global"
	"imgGo-Server/model"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// FileDelete 处理删除文件
func FileDelete(w http.ResponseWriter, r *http.Request) {
	conf := global.Conf
	MyDebug := global.MyDebug

	MyDebug("--call handleDelete--")

	if r.Method != "POST" {
		IndexView(w, r)
		return
	}

	resData := model.DeleteResModel{BaseResModel: model.BaseResModel{Code: 200, Msg: "ok"}}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		resData.Code = 400
		resData.Msg = fmt.Sprintf("read body err, %v\n", err)
		return
	}
	MyDebug("request json ->", string(body))

	var reqData model.DeleteReqData
	if err = json.Unmarshal(body, &reqData); err != nil {
		resData.Code = 400
		resData.Msg = fmt.Sprintf("Unmarshal err, %v\n", err)
		return
	}

	for i, p := range reqData.Paths {
		filePath, fileName := parsePath(p)
		MyDebug(fmt.Sprintf("Deleting... %d/%d", i+1, len(reqData.Paths)), "->", p)
		// 物理删除 TODO 考虑是否增加软删除（移动端一个目录，定时任务清理 或 数据库标识，定时任务物理删除）
		err := os.Remove(path.Join(conf.Path, filePath, fileName))
		if err != nil {
			resData.Code = 400
			resData.Msg = "Delete fail, err: " + path.Join(conf.Path, filePath, fileName)
			MyDebug(resData.Msg)
			break
		}
		// 更新数据库：文件状态
		db := dao.GetDao()
		db.Model(&model.FileDbModel{}).
			Where(fmt.Sprintf("file_path = '%s' and file_name = '%s'", filePath, fileName)).
			Update("status", 0)
		//Updates(FileDbModel{Status: 0, UpdatedAt: time.Now()})
		resData.Data = append(resData.Data, p)
	}

	_data, _ := json.Marshal(resData)

	w.WriteHeader(resData.Code)
	_, err = w.Write(_data)
	if err != nil {
		MyDebug(err)
	}
	return
}

// 解析传入的删除路径，返回标准日期路径、文件名
// 例如：2023/08		21_16-02-33.388.png.webp
// 或：  2022/09/26	1642834046009692757.png
func parsePath(p string) (string, string) {
	p = strings.Replace(p, global.Conf.Url, "", -1)
	_list := strings.Split(p, "/")
	return strings.Join(_list[:len(_list)-1], "/"), _list[len(_list)-1]
}
