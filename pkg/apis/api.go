package apis

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k-avy/CheckOut/pkg/db"
	"github.com/k-avy/CheckOut/pkg/model"
)

func Router() *gin.Engine {
	r := gin.Default()
	authorized := r.Group("/api/v1")
	authorized.GET("/orders", GetallOrders)
	authorized.GET("/orders/:order_id", GetOrder)
	authorized.PUT("/orders/:order_id", UpdateOrder)
	authorized.POST("/orders/:order_id", CreateOrder)
	authorized.DELETE("/orders/:order_id", DeleteOrder)

	return r

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

	order := model.Order(input)

	result := db.DB.
		Where(&model.Order{OrderId: order.OrderId}).
		FirstOrCreate(&order)
	if result.RowsAffected == 0 {
		er := model.OrderError{
			Status:    http.StatusConflict,
			Message:   "can not create order with same order_id",
			Error:     fmt.Sprintf("Error: Order with Order Id '%d' already exists", input.OrderId),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusConflict, er)
		return
	}
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
	res := model.OrdersResponse{
		Data:   []model.Order{order},
		Status: http.StatusCreated,
	}

	c.JSON(http.StatusCreated, res)
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

	res := model.OrdersResponse{
		Data:   []model.Order{order},
		Status: http.StatusOK,
	}

	c.JSON(http.StatusOK, res)
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

	res := model.OrdersResponse{
		Data:   []model.Order{order},
		Status: http.StatusOK,
	}

	c.JSON(http.StatusOK, res)
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

	if order.OrderId!=input.OrderId{
		er := model.OrderError{
			Status:    http.StatusBadRequest,
			Message:   "invalid json body",
			Error:     "id is not matched",
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}

	db.DB.Save(&input)
}

func Middleware(){
	
}