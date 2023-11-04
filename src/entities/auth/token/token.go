package token

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"zgrabi-mjesto.hr/backend/src/config"
)

func GenerateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = strconv.FormatUint(uint64(userId), 10)
	claims["exp"] = time.Now().Add(7 * 24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Config.ApiSecret))

}

func TokenValid(c *gin.Context) error {
	tokenString, err := ExtractToken(c)
	if err != nil {
		return err
	}

	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config.ApiSecret), nil
	})
	if err != nil {
		return err
	}

	return nil
}

func ExtractToken(c *gin.Context) (token string, err error) {
	// Token in query string
	{
		token = c.Query("token")
		if token != "" {
			return
		}
	}

	// Token in authorization header
	{
		header := c.Request.Header.Get("Authorization")
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			err = fmt.Errorf("invalid authorization header. expected format: Authorization: Bearer {token}")
			return
		}

		if strings.ToLower(headerParts[0]) != "bearer" {
			err = fmt.Errorf("invalid authorization header. expected format: Authorization: Bearer {token}")
			return
		}

		token = headerParts[1]

		if token != "" {
			return
		}
	}

	return "", fmt.Errorf("token not found")
}

func ExtractUserIdFromToken(c *gin.Context) (uint, error) {
	tokenString, err := ExtractToken(c)
	if err != nil {
		return 0, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config.ApiSecret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uidString, ok := claims["user_id"].(string)
		if !ok {
			return 0, fmt.Errorf("invalid token")
		}

		uid, err := strconv.ParseUint(uidString, 10, 64)
		if err != nil {
			return 0, err
		}

		return uint(uid), nil
	}
	return 0, nil
}
