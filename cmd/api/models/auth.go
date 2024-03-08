package models

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	bearerPrefix = "Bearer "
)

func GetToken(c *gin.Context) (*jwt.Token, error) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			return nil, errors.New("no token data to receive")
		}
		tokenStr = authHeader[len(bearerPrefix):]
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		secretKey := os.Getenv("secret")
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token cannot be parsed or isn't valid")
	}
	return token, nil
}
