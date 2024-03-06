package model

import "github.com/dgrijalva/jwt-go"

// PaginationModel 分页模型
type PaginationModel struct {
	Index     int
	Size      int
	ItemCount int64
	PageCount int
}

// BaseResModel 基础返回结构体
type BaseResModel struct {
	Code int    `json:"code" yaml:"Code"`
	Msg  string `json:"msg" yaml:"Msg"`
}

// UploadResModel 上传json响应结构体
type UploadResModel struct {
	BaseResModel
	Data []UploadResData `json:"data" yaml:"Data"`
}

// UploadResData 上传数据结构体
type UploadResData struct {
	Name string `json:"name" yaml:"Name"`
	Url  string `json:"url" yaml:"Url"`
}

// DeleteResModel 删除json响应结构体
type DeleteResModel struct {
	BaseResModel
	Data []string `json:"data" yaml:"Data"`
}

// DeleteReqData 删除请求数据结构体
type DeleteReqData struct {
	Paths []string `json:"paths" yaml:"Paths"`
}

// UserAuthResModel 用户认证响应结构体
type UserAuthResModel struct {
	BaseResModel
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type JwtClaims struct {
	ID       uint   `json:"id"`
	Token    string `json:"token"`
	Username string `json:"username"`
	jwt.StandardClaims
}
