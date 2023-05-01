package authbundle

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

	key := []byte(core.Config.Server.JWTPrivateKey)
	return token.SignedString(key)
}

func ValidateJWT(c *gin.Context) error {
	token, err := getJWTokenFromBearerString(c.Request.Header.Get("Authorization"))

	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}

	return errors.New("invalid token provided")
}

func getJWTokenFromBearerString(bearerString string) (*jwt.Token, error) {
	if bearerString == "" {
		return nil, errors.New("empty bearer token")
	}
	splitToken := strings.Split(bearerString, " ")

	var authString string
	if len(splitToken) == 2 {
		authString = splitToken[1]
	}

	token, err := jwt.Parse(authString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key := []byte(core.Config.Server.JWTPrivateKey)

		return key, nil
	})

	return token, err
}

func getBearerTokenFromString(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, "")

	if len(splitToken) == 2 {
		return splitToken[1]
	}

	return ""
}
