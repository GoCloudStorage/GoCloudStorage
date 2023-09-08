package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const issuer = "duryun"
const expireTime = 3000 * time.Second

type UploadToken struct {
	Hash    string
	PartNum int
	Size    int
	jwt.StandardClaims
}

func GenerateToken(hash string, partNum int, size int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(expireTime)
	issuer := issuer
	claims := UploadToken{
		Hash:    hash,
		PartNum: partNum,
		Size:    size,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("golang"))
	return token, err
}

func ParseToken(token string) (*UploadToken, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UploadToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("golang"), nil
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
