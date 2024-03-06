package handle

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"imgGo-Server/dao"
	"imgGo-Server/global"
	"imgGo-Server/model"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type FilerunParams struct {
	Ids   []string `json:"ids[]" yaml:"ids"`
	Paths []string `json:"paths[]" yaml:"paths"`
	Csrf  string   `json:"csrf" yaml:"csrf"`
}

// Filerun 处理Filerun 请求
func Filerun(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery == "module=fileman&section=do&page=delete" {
		handleFilerunDeleteFile(w, r)
	} else if r.URL.RawQuery == "module=trash&section=ajax&page=restore" {
		handleFilerunRestoreFile(w, r)
	} else {
		global.MyDebug("无效的请求参数 ->", r.URL.RawQuery)
	}
}

// handleFilerunDeleteFile 处理Filerun 删除文件
func handleFilerunDeleteFile(w http.ResponseWriter, r *http.Request) {
	/*
		1. 检查是否imgGo 文件；
		2. 更新imgGo 文件状态;
		3. 调用filerun接口。
	*/
	conf := global.Conf
	MyDebug := global.MyDebug
	db := dao.GetDao()

	if conf.IsDebug {
		MyDebug("--call handleFilerunDeleteFile--")
	}

	if r.Method != "POST" {
		SyncViewLogsView(w, r)
		return
	}

	params := FilerunParams{
		Paths: r.Form["paths[]"],
		Csrf:  r.PostFormValue("csrf"),
	}

	MyDebug("Request params ->", params)

	for _, path := range params.Paths {
		if !strings.HasPrefix(path, "/ROOT/HOME/imgGo_upload/") {
			MyDebug("跳过非imgGo 文件 ->", path)
			continue
		}

		// path like: /ROOT/HOME/imgGo_upload/2023/06/25_03-32-17.972.xlsx
		recordItem := model.FileDbModel{}
		filePath := path[24:31]
		fileName := path[32:]
		MyDebug("处理 ->", filePath, fileName)
		query := fmt.Sprintf("file_path='%s' and file_name='%s'", filePath, fileName)
		db.Model(model.FileDbModel{}).Where(query).First(&recordItem)
		recordItem.Status = 0
		if recordItem.Id != 0 {
			db.Save(&recordItem)
			MyDebug("更新状态 -> Status =", recordItem.Status)
		}
	}

	requestFilerun(w, r)
}

// handleFilerunRestoreFile 处理Filerun 恢复已删除的文件
func handleFilerunRestoreFile(w http.ResponseWriter, r *http.Request) {
	/*
		1. filerun 根据id 查找文件路径；
		2. 检查是否imgGo 文件;
		3. 更新imgGo 文件状态；
		4. 调用filerun接口。
	*/
	MyDebug := global.MyDebug
	fDB := dao.GetFilerunDao()
	db := dao.GetDao()

	MyDebug("--call handleFilerunRestoreFile--")
	if r.Method != "POST" {
		SyncViewLogsView(w, r)
		return
	}

	params := FilerunParams{
		// filerun 文件id
		Ids:  r.Form["ids[]"],
		Csrf: r.PostFormValue("csrf"),
	}

	MyDebug("params ->", params)

	var trashes []model.FilerunTrashDbModel
	query := fmt.Sprintf("uid=1 and id in (%s)", strings.Join(params.Ids, ", "))
	fDB.Model(model.FilerunTrashDbModel{}).Where(query).Find(&trashes)

	for _, t := range trashes {
		if strings.HasPrefix(t.RelativePath, "/ROOT/HOME/imgGo_upload/") {
			// path like: /ROOT/HOME/imgGo_upload/2023/06/25_03-32-17.972.xlsx
			recordItem := model.FileDbModel{}
			query := fmt.Sprintf("file_path='%s' and file_name='%s'", t.RelativePath[24:31], t.RelativePath[32:])
			db.Model(model.FileDbModel{}).Where(query).First(&recordItem)
			recordItem.Status = 1
			if recordItem.Id != 0 {
				db.Save(&recordItem)
			}
		}
		// filerun 删除trash
		// fDB.Delete(&t)
	}

	requestFilerun(w, r)
}

// 转发请求给filerun 接口
func requestFilerun(w http.ResponseWriter, r *http.Request) {
	MyDebug := global.MyDebug
	cfg := global.Conf

	MyDebug("--call requestFilerun--")
	// 忽略TLS证书错误
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 创建HTTP客户端
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	reqUrl := cfg.FilerunUrl + "?" + r.URL.RawQuery

	MyDebug("r.UploadMethod ->", reqUrl)

	formData := io.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))

	// 构造一个新的请求，将客户端请求转发到filerun
	req, err := http.NewRequest(r.Method, reqUrl, formData)
	if err != nil {
		log.Fatal(err)
	}

	MyDebug("\n\n复制所有请求头到新的请求中")

	// 复制所有请求头到新的请求中
	for key, value := range r.Header {
		MyDebug(key, value[0])
		req.Header.Set(key, value[0])
	}

	// 发送新请求到filerun
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		MyDebug("Error reading doUploadResponse:", err)
		return
	}

	MyDebug("\n\n将filerun 返回的结果返回给客户端")
	MyDebug("StatusCode ->", resp.StatusCode)
	MyDebug("Body ->", string(body))

	// 将filerun 返回的结果返回给客户端
	for key, value := range resp.Header {
		MyDebug(key, value[0])
		w.Header().Set(key, value[0])
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
