package handle

import (
	"github.com/dgrijalva/jwt-go"
	"imgGo-Server/global"
	"imgGo-Server/model"
	"time"
)

// GenerateToken 生成JWT Token
func GenerateToken(user *model.UserDbModel) (*model.JwtClaims, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(7 * 24 * time.Hour) // Token 7 天过期
	claims := model.JwtClaims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	signedString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(global.Conf.JWTSign))
	claims.Token = signedString
	return &claims, err

}

// ParseToken 解析JWT Token
func ParseToken(token string) (*model.JwtClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &model.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Conf.JWTSign), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*model.JwtClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
