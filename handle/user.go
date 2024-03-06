package handle

import (
	"encoding/json"
	"imgGo-Server/dao"
	"imgGo-Server/global"
	"imgGo-Server/model"
	"imgGo-Server/util"
	"net/http"
	"strings"
)

// UserReg 处理用户注册
func UserReg(w http.ResponseWriter, r *http.Request) {
	conf := global.Conf
	MyDebug := global.MyDebug

	MyDebug("--call UserReg--")

	if !conf.UserConfig.AllowRegister || r.Method != "POST" {
		util.ResponseComm(w, 400, "无效的请求！")
		return
	}

	username := strings.ToLower(r.PostFormValue("username"))
	password := r.PostFormValue("password")
	password = util.Md5Crypt(password, conf.MD5Salt, username)

	user := model.UserDbModel{Username: username, Password: password}

	orm := dao.GetDao()

	result := orm.Create(&user) // 通过数据的指针来创建
	if result.RowsAffected == 1 {
		util.ResponseComm(w, 200, "注册成功！")
		_data, _ := json.Marshal(user)
		MyDebug("user =", string(_data))
	} else {
		errMsg := "注册失败："
		//if errors.Is(result.Error, gorm.ErrDuplicatedKey) { // 错误内容和ErrDuplicatedKey定义的不一样，无法判断
		if strings.HasPrefix(result.Error.Error(), "Error 1062 (23000)") {
			errMsg += "该用户已存在！"
		} else {
			errMsg += "未知错误！"
		}
		util.ResponseComm(w, 400, errMsg)
	}
}

// UserLogin 处理用户登录
func UserLogin(w http.ResponseWriter, r *http.Request) {
	conf := global.Conf
	MyDebug := global.MyDebug

	MyDebug("--call UserLogin--")

	if !conf.UserConfig.AllowLogin || r.Method != "POST" {
		util.ResponseComm(w, 400, "无效的请求！")
		return
	}

	username := strings.ToLower(r.PostFormValue("username"))
	password := r.PostFormValue("password")

	var user model.UserDbModel

	orm := dao.GetDao()
	result := orm.Where("username = ?", username).First(&user)

	if result.RowsAffected == 0 {
		util.ResponseComm(w, 404, "登录失败：用户名不存在！")
		return
	}

	if user.Password != util.Md5Crypt(password, conf.MD5Salt, username) {
		util.ResponseComm(w, 401, "登录失败：密码错误！")
		return
	}

	if user.IsDelete {
		util.ResponseComm(w, 401, "登录失败：该用户已停用！")
		return
	}

	token, err := GenerateToken(&user)
	if err != nil {
		util.ResponseComm(w, 500, "未知错误，请联系管理员！")
		return
	}

	//w.WriteHeader(http.StatusCreated)
	//w.Header().Set("Content-Type", "application/json")

	util.ResponseToken(w, &model.UserAuthResModel{
		BaseResModel: model.BaseResModel{Code: 200, Msg: "ok"},
		Token:        token.Token,
		ExpiresAt:    token.ExpiresAt,
	})

	_ = dao.RecordingNewToken(token)
}
