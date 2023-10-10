package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	secretKey = "fioahgoiankzmcoiwhuooiashniofuaohnrfe"
)

type DownloadToken struct {
	StorageID uint64
	Filename  string
	Ext       string
	jwt.StandardClaims
}

func GenerateDownloadToken(storageID uint64, filename, ext string, expire time.Duration) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(expire)
	logrus.Info(expireTime)
	issuer := issuer
	claims := DownloadToken{
		StorageID: storageID,
		Filename:  filename,
		Ext:       ext,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	return token, err
}

func ParseDownloadToken(token string) (*DownloadToken, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &DownloadToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*DownloadToken); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
