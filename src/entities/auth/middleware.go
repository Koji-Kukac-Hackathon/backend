package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"zgrabi-mjesto.hr/backend/src/entities/auth/roles"
	"zgrabi-mjesto.hr/backend/src/entities/auth/token"
	"zgrabi-mjesto.hr/backend/src/providers/database"
	"zgrabi-mjesto.hr/backend/src/server/response"
)

func respondUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, response.Error("unauthorized"))
	c.Abort()
}

const CtxUserKey = "user"

func extractUserFromToken(c *gin.Context) (*User, error) {
	err := token.TokenValid(c)
	if err != nil {
		return nil, err
	}

	userId, err := token.ExtractUserIdFromToken(c)
	if err != nil {
		return nil, err
	}

	var user *User
	err = database.DatabaseProvider().Client().Model(&User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func RequireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := extractUserFromToken(c)
		if err != nil {
			fmt.Printf("auth error: %s\n", err)
			respondUnauthorized(c)
			return
		}

		c.Set(CtxUserKey, user)

		c.Next()
	}
}

func RequireAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *User
		user, exists := GetUserFromContext(c)
		if !exists {
			tokenUser, err := extractUserFromToken(c)

			if err == nil {
				user = tokenUser
			}
		}

		if user == nil {
			respondUnauthorized(c)
			return
		}

		if user.Role != roles.RoleAdmin {
			respondUnauthorized(c)
			return
		}

		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) (*User, bool) {
	ctxUser, exists := c.Get(CtxUserKey)

	if !exists {
		return nil, false
	}

	user, ok := ctxUser.(*User)

	return user, ok
}
