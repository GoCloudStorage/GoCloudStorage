package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const issuer = "duryun"

type UploadToken struct {
	Key string
	jwt.StandardClaims
}

func GenerateUploadToken(key string, expire time.Duration) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(expire)
	issuer := issuer
	claims := UploadToken{
		Key: key,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	return token, err
}

func ParseUploadToken(token string) (*UploadToken, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UploadToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UploadToken); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
