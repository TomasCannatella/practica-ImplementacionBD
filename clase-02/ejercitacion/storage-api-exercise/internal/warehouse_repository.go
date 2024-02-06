package internal

import (
	"errors"
)

var (
	ErrWarehouseNotFound      = errors.New("repository: warehouse not found")
	ErrWarehouseAlreadyExists = errors.New("repository: warehouse already exists")
)

type WarehouseRepository interface {
	// GetAll returns all warehouses
	GetAll() (w []Warehouse, err error)
	// GetOne returns a warehouse by id
	GetOne(id int) (w Warehouse, err error)
	// Store saves a warehouse
	Store(w *Warehouse) (err error)
	// ReportProducts returns a report of products by warehouse
	ReportProducts(id int) (rp []ReportProduct, err error)
}
