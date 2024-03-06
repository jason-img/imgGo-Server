package util

import (
	"encoding/json"
	"imgGo-Server/global"
	"imgGo-Server/model"
	"net/http"
)

// ResponseUpload 返回客户端
func ResponseUpload(w http.ResponseWriter, resData *model.UploadResModel) {
	msg, _ := json.Marshal(resData)
	//w.Header().Set("content-type", "application/json")
	w.WriteHeader(resData.Code)
	_, err := w.Write(msg)
	if err != nil {
		global.MyDebug(err)
	}
}

// ResponseComm 返回客户端
func ResponseComm(w http.ResponseWriter, code int, msg string) {
	resData := model.UploadResModel{}
	resData.Code = code
	resData.Msg = msg
	data, _ := json.Marshal(resData)
	//w.Header().Set("content-type", "application/json")
	w.WriteHeader(resData.Code)
	_, err := w.Write(data)
	if err != nil {
		global.MyDebug(err)
	}

}

// ResponseToken 返回token
func ResponseToken(w http.ResponseWriter, resData *model.UserAuthResModel) {
	msg, _ := json.Marshal(resData)
	//w.Header().Set("content-type", "application/json")
	w.WriteHeader(resData.Code)
	_, err := w.Write(msg)
	if err != nil {
		global.MyDebug(err)
	}
}
