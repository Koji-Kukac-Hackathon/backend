package auth

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"zgrabi-mjesto.hr/backend/src/entities/auth/token"
	"zgrabi-mjesto.hr/backend/src/providers/database"
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

func (controller) ListUsers(c *gin.Context) {
	type QueryParams struct {
		Limit   uint
		OrderBy string
		Order   string
		After   *uint
		Filter  *string
	}

	queryParams := &QueryParams{
		Limit:   50,
		OrderBy: "created_at",
		Order:   "desc",
		After:   nil,
		Filter:  nil,
	}

	var allowedOrders = map[string]bool{
		"asc":  true,
		"desc": true,
	}
	var allowedColumns = map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"role":       true,
		"created_at": true,
	}

	{
		limitQuery := c.Query("limit")
		limit64, err := strconv.ParseUint(limitQuery, 10, 32)
		if err == nil {
			queryParams.Limit = max(1, min(uint(limit64), 50))
		}
	}

	{
		orderByQuery := c.Query("order_by")
		if _, ok := allowedColumns[orderByQuery]; ok {
			queryParams.OrderBy = orderByQuery
		}
	}

	{
		orderQuery := c.Query("order")
		if _, ok := allowedOrders[orderQuery]; ok {
			queryParams.Order = orderQuery
		}
	}

	{
		filterQuery := c.Query("filter")
		if filterQuery != "" {
			queryParams.Filter = &filterQuery
		}
	}

	{
		afterQuery := c.Query("after")
		after64, err := strconv.ParseUint(afterQuery, 10, 64)
		if err == nil {
			after := uint(after64)
			queryParams.After = &after
		}
	}

	db := database.DatabaseProvider().Client()

	var users []User
	{
		res := db.Model(&User{}).Limit(int(queryParams.Limit)).Order(fmt.Sprintf("%s %s", queryParams.OrderBy, queryParams.Order))

		if queryParams.Filter != nil {
			res = res.Where("1 = 0")
			res = res.Or("name LIKE ?", fmt.Sprintf("%%%s%%", *queryParams.Filter))
			res = res.Or("email LIKE ?", fmt.Sprintf("%%%s%%", *queryParams.Filter))
			res = res.Or("role LIKE ?", fmt.Sprintf("%%%s%%", *queryParams.Filter))
		}

		if queryParams.After != nil {
			res = res.Where("id > ?", *queryParams.After)
		}

		if err := res.Find(&users).Error; err != nil {
			c.JSON(http.StatusBadRequest, response.Error(err.Error()))
			return
		}
	}

	var totalUsers int64
	{
		res := db.Model(&User{}).Count(&totalUsers)
		if err := res.Error; err != nil {
			c.JSON(http.StatusBadRequest, response.Error(err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, response.Success("List of users", gin.H{
		"items":      users,
		"filter":     queryParams,
		"totalItems": totalUsers,
	}))
}

type UpdateUserData struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Role     string `json:"role" form:"role"`
	Password string `json:"password" form:"password"`
}

func (controller) UpdateUser(c *gin.Context) {
	userIdQuery := c.Param("id")
	userId64, err := strconv.ParseUint(userIdQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	var updateUserData UpdateUserData
	if err := c.BindJSON(&updateUserData); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	if err := Service.EditUser(uint(userId64), &updateUserData); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Successfully updated user"))
}

func (controller) DeleteUser(c *gin.Context) {
	userIdQuery := c.Param("id")
	userId64, err := strconv.ParseUint(userIdQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	if err := Service.DeleteUser(uint(userId64)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Successfully deleted user"))
}
