package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"zgrabi-mjesto.hr/backend/src/entities/auth/token"
	"zgrabi-mjesto.hr/backend/src/server/response"
)

type controller struct{}

var Controller controller = controller{}

func (controller) GetUser(c *gin.Context) {
	user, exists := GetUserFromContext(c)

	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error("Unauthorized"))
		return
	}

	c.JSON(http.StatusOK, response.Success("User info", user))
}

type LoginData struct {
	Email    *string `json:"email" form:"email"`
	Password *string `json:"password" form:"password"`
	Provider string  `json:"provider" form:"provider"`
}

func (controller) Login(c *gin.Context) {
	var loginData LoginData
	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	user, err := Service.Login(&loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	jwtToken, err := token.GenerateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	fmt.Printf("Login: %+v\n", loginData)
	c.JSON(http.StatusOK, response.Success("Successfully authenticated user", gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"token": jwtToken,
	}))
}

type RegisterData struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (controller) Register(c *gin.Context) {
	var registerData RegisterData
	if err := c.BindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	err := Service.Register(&registerData)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	fmt.Printf("Register: %+v\n", registerData)
	c.JSON(http.StatusOK, response.Success("Successfully registered user"))
}
