package middleware

import (
	"imgGo-Server/global"
	"imgGo-Server/handle"
	"imgGo-Server/util"
	"net/http"
)

// Auth 用户认证
func Auth(next http.Handler) http.HandlerFunc {
	MyDebug := global.MyDebug
	return func(w http.ResponseWriter, r *http.Request) {
		// 在处理请求之前执行的逻辑
		MyDebug("Before Auth ->", r.URL)

		// 统一设置content-type
		w.Header().Set("content-type", "application/json; charset=utf-8")

		token := r.Header.Get("authorization")
		if token == "" {
			util.ResponseComm(w, 401, "Token 不能为空！")
			return
		}
		if _, err := handle.ParseToken(token); err != nil {
			MyDebug("解析Token 失败", err)
			util.ResponseComm(w, 401, "无效的Token")
			return
		}

		// 调用下一个处理器
		next.ServeHTTP(w, r)
		// 在处理请求之后执行的逻辑
		MyDebug("After Auth ->", r.URL)
	}
}

// UserPwdCheck 验证用户名或密码是否符合规范
func UserPwdCheck(next http.Handler) http.HandlerFunc {
	conf := global.Conf.UserConfig
	MyDebug := global.MyDebug
	return func(w http.ResponseWriter, r *http.Request) {
		// 在处理请求之前执行的逻辑
		MyDebug("Before UserPwdCheck ->", r.URL)

		// 统一设置content-type
		w.Header().Set("content-type", "application/json; charset=utf-8")

		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if username == "" || password == "" {
			util.ResponseComm(w, 400, "用户名、密码不能为空！")
			return
		}

		if !util.RegexpValid(username, conf.UsernameRegex) {
			util.ResponseComm(w, 400, "无效的用户名，请输入3~16位字母或数字！")
			return
		}

		if !util.RegexpValid(password, conf.PwdRegex) {
			util.ResponseComm(w, 400, "无效的密码，请输入8~32位包含大小写字母和数字的有效密码！")
			return
		}

		// 调用下一个处理器
		next.ServeHTTP(w, r)

		// 在处理请求之后执行的逻辑
		MyDebug("After UserPwdCheck:", r.URL)
	}
}
