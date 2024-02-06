package repository

import (
	"app/internal"
	"database/sql"
	"errors"
)

// NewProductsMySQL returns a new instance of ProductsMySQL
func NewProductsMySQL(db *sql.DB) *ProductsMySQL {
	return &ProductsMySQL{
		db: db,
	}
}

// ProductsMySQL is a struct that represents a product repository
type ProductsMySQL struct {
	// db is the database connection
	db *sql.DB
}

// GetAll returns all products
func (r *ProductsMySQL) GetAll() (products []internal.Product, err error) {
	query := "SELECT `id`, `name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`, `id_warehouse` FROM `products`"
	
	row, err := r.db.Query(query)
	if err != nil {
		if errors.Is(err,sql.ErrNoRows) {
			err = internal.ErrProductNotFound
			return
		}
		return
	}
	defer row.Close()

	var p internal.Product
	for row.Next(){
		err = row.Scan(&p.ID,&p.Name,&p.Quantity,&p.CodeValue,&p.IsPublished,&p.Expiration,&p.Price,&p.WarehouseId)
		if err != nil {
			return
		}
		products = append(products,p)
	}

	return
}

// GetOne returns a product by id
func (r *ProductsMySQL) GetOne(id int) (p internal.Product, err error) {
	// execute the query
	row := r.db.QueryRow(
		"SELECT `id`, `name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`, `id_warehouse` "+
			"FROM `products` WHERE `id` = ?",
		id,
	)
	if err = row.Err(); err != nil {
		return
	}

	// scan the row into the product
	err = row.Scan(&p.ID, &p.Name, &p.Quantity, &p.CodeValue, &p.IsPublished, &p.Expiration, &p.Price, &p.WarehouseId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrProductNotFound
		}
		return
	}

	return
}

// Store stores a product
func (r *ProductsMySQL) Store(p *internal.Product) (err error) {
	// execute the query
	result, err := r.db.Exec(
		"INSERT INTO `products` (`name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`, `id_warehouse`) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?)",
		p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price, p.WarehouseId,
	)
	if err != nil {
		return
	}

	// get the last inserted id
	id, err := result.LastInsertId()
	if err != nil {
		return
	}
	p.ID = int(id)

	return
}

// Update updates a product
func (r *ProductsMySQL) Update(p *internal.Product) (err error) {
	// execute the query
	_, err = r.db.Exec(
		"UPDATE `products` SET `name` = ?, `quantity` = ?, `code_value` = ?, `is_published` = ?, `expiration` = ?, `price` = ? , `id_warehouse` = ?"+
			"WHERE `id` = ?",
		p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price, p.WarehouseId, p.ID,
	)
	if err != nil {
		return
	}

	return
}

// Delete deletes a product by id
func (r *ProductsMySQL) Delete(id int) (err error) {
	// execute the query
	_, err = r.db.Exec(
		"DELETE FROM `products` WHERE `id` = ?",
		id,
	)
	if err != nil {
		return
	}

	return
}
