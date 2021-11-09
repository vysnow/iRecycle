package common

import (
	"time"

	"com.mego.first/megofirst/model"
	"github.com/golang-jwt/jwt"
)

var jwtToken = []byte("my_jw7_5ecret_t0k3n")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func IssueToken(user model.User) (string, error) {
	expiresAt := time.Now().Add(2 * time.Hour).Unix()
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Mego First",
			Subject:   "Auth Token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtToken)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtToken, nil
	})

	return token, claims, err
}
