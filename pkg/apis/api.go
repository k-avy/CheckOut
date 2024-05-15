package apis

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k-avy/CheckOut/pkg/db"
	"github.com/k-avy/CheckOut/pkg/model"
)

func Router(rl int64) *gin.Engine {
	r := gin.Default()
	r.Use(RateLimiter(rl))
	authorized := r.Group("/api")
	authorized.Use(AuthMiddleware())
	authorized.GET("/orders", GetallOrders)
	authorized.GET("/orders/:order_id", GetOrder)
	authorized.PUT("/orders/:order_id", UpdateOrder)
	authorized.POST("/orders", CreateOrder)
	authorized.DELETE("/orders/:order_id", DeleteOrder)

	r.POST("/register", Register)

	return r

}

func Register(c *gin.Context) {
	var input model.User

	if err := c.ShouldBindJSON(&input); err != nil {
		er := model.OrderError{
			Status:    http.StatusBadRequest,
			Message:   "wrong Json body",
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}

	result := db.DB.FirstOrCreate(&input, model.User{Username: input.Username})
	if result.Error != nil {
		er := model.OrderError{
			Status:    http.StatusBadRequest,
			Message:   "Failed to register",
			Error:     result.Error.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}

	if result.RowsAffected == 0 {
		er := model.OrderError{
			Status:    http.StatusBadRequest,
			Message:   "username exists",
			Error:     "Error: duplicate registration",
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			er := model.OrderError{
				Status:    http.StatusUnauthorized,
				Message:   "not authorized",
				Error:     "no credentials provided",
				Path:      c.Request.URL.Path,
				Timestamp: time.Now().Format(time.RFC3339),
			}
			c.JSON(http.StatusUnauthorized, er)
			c.Abort()
			return
		}

		authString, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))

		if err != nil {
			er := model.OrderError{
				Status:    http.StatusUnauthorized,
				Message:   "not authorized, check credentials",
				Error:     err.Error(),
				Path:      c.Request.URL.Path,
				Timestamp: time.Now().Format(time.RFC3339),
			}
			c.JSON(http.StatusUnauthorized, er)
			c.Abort()
			return
		}

		auth := strings.SplitN(strings.TrimSpace(string(authString)), ":", 2)
		if len(auth) != 2 {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		username := auth[0]
		password := auth[1]

		var user model.User

		errd := db.DB.Where("username = ?", username).First(&user).Error

		if errd != nil {
			er := model.OrderError{
				Status:    http.StatusUnauthorized,
				Message:   "user not found",
				Error:     errd.Error(),
				Path:      c.Request.URL.Path,
				Timestamp: time.Now().Format(time.RFC3339),
			}
			c.JSON(http.StatusUnauthorized, er)
			c.Abort()
			return
		}

		if strings.Compare(user.Password, password) != 0 {
			er := model.OrderError{
				Status:    http.StatusUnauthorized,
				Message:   "not authorized",
				Error:     "wrong password",
				Path:      c.Request.URL.Path,
				Timestamp: time.Now().Format(time.RFC3339),
			}
			c.JSON(http.StatusUnauthorized, er)
			c.Abort()
			return
		}

		c.Next()
	}
}

func RateLimiter(rl int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !fixedWindowRateLimiter("global_rate_limiter", rl, time.Minute) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		}
		c.Next()
	}
}

func fixedWindowRateLimiter(key string, rl int64, r time.Duration) bool {
	keyWindow := fmt.Sprintf("%s:%d", key, time.Now().Unix()/int64(r.Seconds()))
	count, err := db.RED.Incr(context.Background(), keyWindow).Result()
	if err != nil {
		panic(err)
	}
	if count == 1 {
		if err := db.RED.Expire(context.Background(), keyWindow, r).Err(); err != nil {
			panic(err)
		}
	}
	return count <= rl
}

func GetallOrders(c *gin.Context) {
	var orders []model.Order
	query := db.DB.Model(&model.Order{})
	err := query.Find(&orders).Error
	if err != nil {
		er := model.OrderError{
			Status:    http.StatusBadRequest,
			Message:   "Orders not found",
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}
	res := model.OrdersResponse{
		Data:   orders,
		Status: http.StatusOK,
	}

	c.JSON(http.StatusOK, res)

}

func CreateOrder(c *gin.Context) {
	var input model.InputOrder
	if err := c.ShouldBindJSON(&input); err != nil {
		er := model.OrderError{
			Status:    http.StatusBadRequest,
			Message:   "invalid json body",
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}

	order := model.Order{
		OrderDate:   input.OrderDate,
		Customer:    input.Customer,
		ProductName: input.ProductName,
		Quantity:    input.Quantity,
		UnitPrice:   input.UnitPrice,
		Priority:    input.Priority,
	}

	result := db.DB.Create(&order)

	if result.Error != nil {
		er := model.OrderError{
			Status:    http.StatusInternalServerError,
			Message:   "Failed to create Order",
			Error:     result.Error.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusInternalServerError, er)
		return
	}

	c.JSON(http.StatusCreated, order)
}

func GetOrder(c *gin.Context) {
	var order model.Order

	queryParam := c.Param("order_id")
	if queryParam == "" {
		return
	}
	i, _ := strconv.Atoi(queryParam)
	err := db.DB.Where(model.Order{OrderId: i}).First(&order).Error

	if err != nil {
		er := model.OrderError{
			Status:    http.StatusNotFound,
			Message:   fmt.Sprintf("no Order with Order id '%s' exists", queryParam),
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusNotFound, er)
		return
	}

	c.JSON(http.StatusOK, order)
}

func DeleteOrder(c *gin.Context) {
	var order model.Order

	queryParam := c.Param("order_id")
	if queryParam == "" {
		return
	}
	i, _ := strconv.Atoi(queryParam)
	err := db.DB.Where(model.Order{OrderId: i}).Delete(&order).Error

	if err != nil {
		er := model.OrderError{
			Status:    http.StatusNotFound,
			Message:   fmt.Sprintf("no Order with Order id '%s' exists", queryParam),
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusNotFound, er)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}

func UpdateOrder(c *gin.Context) {
	var order model.Order

	queryParam := c.Param("order_id")
	if queryParam == "" {
		return
	}
	i, _ := strconv.Atoi(queryParam)
	err := db.DB.Where(model.Order{OrderId: i}).First(&order).Error

	if err != nil {
		er := model.OrderError{
			Status:    http.StatusNotFound,
			Message:   fmt.Sprintf("no Order with Order id '%s' exists", queryParam),
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusNotFound, er)
		return
	}

	if order.OrderId != i {
		er := model.OrderError{
			Status:    http.StatusBadRequest,
			Message:   "trying to update a different order",
			Error:     "id is not matched",
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}

	var input model.InputOrder
	if err := c.ShouldBindJSON(&input); err != nil {
		er := model.OrderError{
			Status:    http.StatusBadRequest,
			Message:   "invalid json body",
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}

	uorder := model.Order{
		OrderId:     i,
		Customer:    input.Customer,
		ProductName: input.ProductName,
		Quantity:    input.Quantity,
		UnitPrice:   input.UnitPrice,
		OrderDate:   input.OrderDate,
		Priority:    input.Priority,
	}

	if err := db.DB.Save(&uorder).Error; err != nil {
		er := model.OrderError{
			Status:    http.StatusInternalServerError,
			Message:   "could not update",
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order updated"})
}
