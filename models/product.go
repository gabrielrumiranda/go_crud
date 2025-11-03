package models

import "time"

type Product struct {
	ID            int       `json:"id"`
	Name          string    `json:"name" binding:"required"`
	Price         float32   `json:"price" binding:"required,min=0"`
	StockQuantity int       `json:"stock_quantity" binding:"required,min=0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name          string  `json:"name" binding:"required"`
	Price         float32 `json:"price" binding:"required,min=0"`
	StockQuantity int     `json:"stock_quantity" binding:"required,min=0"`
}

type UpdateProductRequest struct {
	Name          string  `json:"name"`
	Price         float32 `json:"price" binding:"omitempty,min=0"`
	StockQuantity int     `json:"stock_quantity" binding:"omitempty,min=0"`
}
