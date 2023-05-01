package authbundle

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	core "github.com/rubenscholle/venus/corebundle"
)

func GenerateJWT(user AuthUser) (string, error) {
	tokenTTL := core.Config.Server.TokenTTL
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"exp":        time.Now().Add(time.Hour * 24 * time.Duration(tokenTTL)).Unix(),
		"authorized": true,
	})

	key := core.Config.Server.JWTPrivateKey
	return token.SignedString(key)
}
