package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dlclark/regexp2"
	"imgGo-Server/global"
	"io"
	"os"
	"strings"
)

// Contains 返回str 是否在s 切片中存在
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// GetLastFileMd5 获取上一次处理的文件MD5码
func GetLastFileMd5() string {
	conf := global.Conf
	global.MyDebug("--call getLastFileMd5--")
	content, err := os.ReadFile(conf.LastMd5File)
	if err != nil {
		_ = fmt.Errorf("打开文件失败，filename=%v, err=%v", conf.LastMd5File, err)
		return ""
	}
	return string(content)
}

// UpdateLastFileMd5 更新上一次处理的文件MD5码
func UpdateLastFileMd5(md5text string) {
	conf := global.Conf
	global.MyDebug("--call updateLastFileMd5--")
	content := []byte(md5text)
	err := os.WriteFile(conf.LastMd5File, content, 0644)
	if err != nil {
		_ = fmt.Errorf("写入文件失败，filename=%v, context=%v, err=%v", conf.LastMd5File, md5text, err)
	}
}

// GetFileMd5 获取指定文件的md5码
func GetFileMd5(filePath string) string {
	global.MyDebug("--call getFileMd5--")
	pFile, err := os.Open(filePath)
	if err != nil {
		_ = fmt.Errorf("打开文件失败，filename=%v, err=%v", filePath, err)
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	_, _ = io.Copy(md5h, pFile)

	return hex.EncodeToString(md5h.Sum(nil))
}

// FileTypeBlocked 根据传入的文件名，返回文件类型是否被屏蔽
func FileTypeBlocked(fileName string) (bool, string) {
	conf := global.Conf
	global.MyDebug("--call fileTypeBlocked--")
	acceptTypes := conf.AcceptFileTypes
	_split := strings.Split(fileName, ".")
	extName := strings.ToLower(_split[len(_split)-1])
	for _, item := range acceptTypes {
		if extName == item {
			return false, extName
		}
	}
	return true, extName
}

// CreateDir 调用os.MkdirAll递归创建文件夹
func CreateDir(filePath string) error {
	global.MyDebug("--call createDir--")
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// IsExist 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// Md5Crypt 给字符串生成md5
// @params plain 需要加密的字符串
// @params salt interface{} 加密的盐
// @return CryptStr 返回md5码
func Md5Crypt(plain string, salts ...interface{}) (CryptStr string) {
	if l := len(salts); l > 0 {
		slice := make([]string, l+1)
		plain = fmt.Sprintf(plain+strings.Join(slice, "%v"), salts...)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(plain)))
}

// RegexpValid 用户名格式验证
func RegexpValid(text string, pattern string) (isValid bool) {
	var err error
	var reg *regexp2.Regexp

	reg, err = regexp2.Compile(pattern, 0)
	if err != nil {
		panic(err)
	}
	isValid, err = reg.MatchString(text)
	if err != nil {
		panic(err)
	}
	return isValid
}
