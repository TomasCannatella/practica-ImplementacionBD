package internal

import "time"

// Product is an struct that represents a product
type Product struct {
	// ID is the unique identifier of the product
	ID int
	// Name is the name of the product
	Name string
	// Quantity is the quantity of the product
	Quantity int
	// CodeValue is the code value of the product
	CodeValue string
	// IsPublished is the published status of the product
	IsPublished bool
	// Expiration is the expiration date of the product
	Expiration time.Time
	// Price is the price of the product
	Price float64
	// WarehouseId is the warehouse id of the product
	WarehouseId int
}
