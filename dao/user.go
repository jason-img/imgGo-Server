package dao

import (
	"imgGo-Server/global"
	"imgGo-Server/model"
	"time"
)

func RecordingNewToken(claims *model.JwtClaims) error {
	MyDebug := global.MyDebug
	MyDebug("----开始记录Token----")
	MyDebug(claims)

	token := model.TokenDbModel{
		Token:     claims.Token,
		IsDelete:  false,
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		Username:  claims.Username,
	}

	result := orm.Create(&token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
