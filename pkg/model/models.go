package model

type Order struct {
	OrderId     int     `json:"order_id" gorm:"primaryKey"`
	Customer    string  `json:"customer"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float32 `json:"unit_price"`
	OrderDate   string  `json:"order_date"`
	Priority    string  `json:"priority"`
}

type OrderError struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Error     string `json:"error"`
	Path      string `json:"path"`
	Timestamp string `json:"timestamp"`
}

type OrdersResponse struct {
	Status int     `json:"status" example:"200"`
	Data   []Order `json:"data"`
}

type InputOrder struct {
	Customer    string  `json:"customer" binding:"required"`
	ProductName string  `json:"product_name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	UnitPrice   float32 `json:"unit_price" binding:"required"`
	OrderDate   string  `json:"order_date" binding:"required"`
	Priority    string  `json:"priority" binding:"required"`
}

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
