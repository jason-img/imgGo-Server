package dao

import (
	"bufio"
	"fmt"
	"gorm.io/gorm"
	"imgGo-Server/global"
	"imgGo-Server/model"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ReadLogFile 按行读取日志文件
func ReadLogFile(filePath string) (int, int) {
	MyDebug := global.MyDebug
	MyDebug("--call ReadLogFile--")
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("打开文件失败，filename=%v, err=%v", filePath, err)
		return 0, 0
	}
	defer file.Close()

	MyDebug("处理日志文件 ->", filePath)

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	var size = stat.Size()
	MyDebug("file size =", size)

	buf := bufio.NewReader(file)
	var lineCount = 0
	var recordCount = 0

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				_ = fmt.Errorf("读取文件失败，filename=%v, err=%v", filePath, err)
				return 0, 0
			}
		}
		line = strings.TrimSpace(line)

		//MyDebug("nginx log line ->", line)

		match, _ := regexp.MatchString(`GET.+HTTP/2.0`, line)
		if match {
			recordCount += MakeViewRecord(line)
			lineCount += 1
		}
	}

	MyDebug("file lines =", lineCount, ", record lines =", recordCount)

	return lineCount, recordCount
}

// MakeViewRecord 创建访问日志 记录
func MakeViewRecord(line string) (createdNum int) {
	MyDebug := global.MyDebug
	MyDebug("--call MakeViewRecord--")
	var item model.ViewLogDbModel
	createdNum = 0

	stepList := strings.Split(line, " - - ")

	//MyDebug("stepList ->", stepList)

	if len(stepList) == 2 {
		item.ClientIp = stepList[0]
		line = stepList[1]
	}

	stringTime := line[1:21]
	loc, _ := time.LoadLocation("Local")
	item.RequestTime, _ = time.ParseInLocation("02/Jan/2006:15:04:05", stringTime, loc)

	line = line[30:]

	stepList = strings.Split(line, " ")

	//MyDebug("stepList ->", stepList)

	item.Method = stepList[0]
	subPathList := strings.Split(stepList[1][1:], "/")
	item.Filename = subPathList[len(subPathList)-1]

	// 匹配imgGo 上传文件
	matched, _ := regexp.Match(`^\d{2}_\d{2}-\d{2}-\d{2}\.\d{3}.+$`, []byte(item.Filename))
	if !matched {
		MyDebug("非imgGo 上传文件 ->", item.Filename, "，跳过...")
		return
	}

	item.HttpVersion = stepList[2][:len(stepList[2])-1]
	item.HttpCode, _ = strconv.Atoi(stepList[3])
	item.Size, _ = strconv.Atoi(stepList[4])
	item.UA = strings.Join(stepList[6:len(stepList)-1], " ")
	item.UA = item.UA[1:]

	MyDebug("item ->", item)

	result := orm.Create(&item)
	if result.Error != nil {
		MyDebug("写入日志记录失败 ->", result.Error)
		return
	}
	createdNum = 1
	if addNum, err := AddViewCount(item.Filename); addNum == 0 || err != nil {
		MyDebug("添加文件访问次数失败 ->", err)
	}
	return
}

// AddViewCount 添加文件访问数
func AddViewCount(fileName string) (int64, error) {
	global.MyDebug("--call AddViewCount--")
	var file model.FileDbModel
	//result := db.Debug().Model(&file).Where("file_name=?", fileName).UpdateColumn("view_count", gorm.Expr("view_count+1"))
	result := orm.Model(&file).Where("file_name=?", fileName).
		UpdateColumn("view_count", gorm.Expr("view_count+1")).
		UpdateColumn("last_view_time", time.Now())
	// 更新文件最近访问时间

	return result.RowsAffected, result.Error
}
