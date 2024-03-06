package handle

import (
	"imgGo-Server/dao"
	"imgGo-Server/global"
	"imgGo-Server/model"
	"imgGo-Server/util"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

// FileUpload 处理文件上传
func FileUpload(w http.ResponseWriter, r *http.Request) {
	conf := global.Conf
	MyDebug := global.MyDebug

	if r.Method != "POST" {
		IndexView(w, r)
		return
	}

	var resData []model.UploadResData
	result := model.UploadResModel{BaseResModel: model.BaseResModel{Code: 200, Msg: "ok"}, Data: resData}

	MyDebug("--call FileUpload--")

	//设置内存大小
	_ = r.ParseMultipartForm(32 << 20)

	//获取上传的文件
	if r.MultipartForm == nil || r.MultipartForm.File == nil || r.MultipartForm.File["files"] == nil {
		util.ResponseComm(w, 404, "文件不能为空")
		return
	}

	// 获取需转换格式
	cFormat := r.PostFormValue("format")
	if cFormat != "" {
		MyDebug("convert format to -->", cFormat)
	}

	var isWebp = false
	var imageWidth int
	var imageHeight int
	var recordItems []model.FileDbModel
	files := r.MultipartForm.File["files"]
	fileCount := strconv.Itoa(len(files))
	datePath := path.Join(time.Now().Format("2006"), time.Now().Format("01"))
	for i, file := range files {
		MyDebug("file.Header -->", file.Header)
		MyDebug("file.Filename -->", file.Filename)
		MyDebug("file.Size -->", file.Size)

		// 检查文件格式是否允许上传
		MyDebug("检查文件格式是否允许上传")
		isBlocked, extName := util.FileTypeBlocked(file.Filename)
		if isBlocked {
			util.ResponseComm(w, 403, "不支持上传此类型的文件: "+extName)
			return
		}
		//打开上传文件
		multipartFile, err := file.Open()
		if multipartFile == nil || err != nil {
			MyDebug(err)
			continue
		}
		defer func(f multipart.File) {
			err = f.Close()
		}(multipartFile)
		MyDebug("打开上传文件 ->", strconv.Itoa(i+1)+"/"+fileCount, file.Filename)

		// 创建上传目录
		// _ = os.Mkdir("./upload", os.ModePerm)

		// 创建上传文件
		// newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + "." + extName
		newFileName := time.Now().Format("02_15-04-05.000") + "." + extName
		MyDebug("创建上传文件 ->", newFileName)
		if !util.IsExist(path.Join(conf.Path, datePath)) {
			if err = util.CreateDir(path.Join(conf.Path, datePath)); err != nil {
				MyDebug(err)
			}
		}
		osFile, err := os.Create(path.Join(conf.Path, datePath, newFileName))
		if osFile == nil || err != nil {
			MyDebug(err)
			continue
		}

		_, _ = io.Copy(osFile, multipartFile)

		// 直接关闭文件，不用defer。否则后面如果需要删除的时候可能还占用着，删不掉
		osFile.Close()

		// 转换为Webp
		if cFormat == "webp" && util.Contains(conf.ConvertToWebpTypes, extName) {
			MyDebug("转换为Webp")
			img := Image{}
			err = img.Open(path.Join(conf.Path, datePath, newFileName), extName)
			if err != nil {
				MyDebug("webp open error -->", err.Error())
				goto DONE
			}

			webpBytes, err := img.ToWebP(conf.WebpQuality)
			if err != nil {
				MyDebug("webp open error -->", err.Error())
				goto DONE

			}

			f, err := os.Create(path.Join(conf.Path, datePath, newFileName+".webp"))
			if err != nil {
				MyDebug("webp open error -->", err.Error())
				goto DONE
			}
			defer f.Close()

			_, err = f.Write(webpBytes)
			if err != nil {
				MyDebug("webp open error -->", err.Error())
				goto DONE
			}

			// 删除原始文件
			if conf.DelOriginFile {
				MyDebug("删除原始文件 ->", newFileName)
				fullPath := path.Join(conf.Path, datePath, newFileName)
				err = os.Remove(fullPath)
				if err != nil {
					MyDebug(err)
					MyDebug("删除原始文件失败 ->", fullPath)
				}
			}

			isWebp = true
			newFileName = newFileName + ".webp"
			imageWidth = img.Width
			imageHeight = img.Height
		}

	DONE:
		// 数据库记录
		recordItem := model.FileDbModel{
			OriginName:   file.Filename,
			FilePath:     datePath,
			FileName:     newFileName,
			FileExtName:  extName,
			IsWebp:       isWebp,
			ImageWidth:   imageWidth,
			ImageHeight:  imageHeight,
			LastViewTime: time.Now(),
			//UpdatedAt:    time.Now(),
		}
		recordItems = append(recordItems, recordItem)

		// 返回数据
		resData = append(resData, model.UploadResData{Name: file.Filename, Url: conf.Url + datePath + "/" + newFileName})
	}

	if len(resData) == 0 {
		result.Code = 404
		result.Msg = "没有可上传的文件"
	} else {
		// 调用协程 创建数据库记录
		go func(records []model.FileDbModel) {
			db := dao.GetDao()
			db.Create(records)
		}(recordItems)
	}

	result.Data = resData
	util.ResponseUpload(w, &result)
	return
}
