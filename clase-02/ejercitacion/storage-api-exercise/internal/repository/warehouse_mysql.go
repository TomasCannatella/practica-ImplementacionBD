package repository

import (
	"app/internal"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

type WarehouseMySQL struct {
	// db is the database connection
	db *sql.DB
}

func NewWarehouseMySQL(db *sql.DB) *WarehouseMySQL {
	return &WarehouseMySQL{
		db: db,
	}
}

func (r *WarehouseMySQL) GetAll() (w []internal.Warehouse, err error) {
	query := "SELECT `id`, `name`, `adress`, `telephone`, `capacity` FROM `warehouses`"
	row, err := r.db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrWarehouseNotFound
			return
		}
		return
	}

	for row.Next() {
		var warehouse internal.Warehouse
		err = row.Scan(&warehouse.Id, &warehouse.Name, &warehouse.Address, &warehouse.Telephone, &warehouse.Capacity)
		if err != nil {
			return
		}
		w = append(w, warehouse)
	}
	return
}

func (r *WarehouseMySQL) GetOne(id int) (w internal.Warehouse, err error) {
	query := "SELECT `id`, `name`, `adress`, `telephone`, `capacity` FROM `warehouses` WHERE `id` = ?"

	row := r.db.QueryRow(query, id)
	if err = row.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrWarehouseNotFound
			return
		}
		return
	}

	err = row.Scan(&w.Id, &w.Name, &w.Address, &w.Telephone, &w.Capacity)
	if err != nil {
		return
	}
	return
}

func (r *WarehouseMySQL) Store(w *internal.Warehouse) (err error) {
	query := "INSERT INTO `warehouses` (`name`, `adress`, `telephone`, `capacity`) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, w.Name, w.Address, w.Telephone, w.Capacity)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErrWarehouseAlreadyExists
				return
			default:
				return
			}
		}
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}
	w.Id = int(id)
	return
}

func (r *WarehouseMySQL) ReportProducts(id int) (rp []internal.ReportProduct, err error) {
	query := "SELECT w.`name`,count(p.id) as `product_count` FROM warehouses w LEFT JOIN products p ON w.id = p.id_warehouse "
	if id > 0 {
		query += "WHERE w.id = ? "
	}
	query += "GROUP BY w.`name`;"

	var result *sql.Rows
	if id > 0 {
		result, err = r.db.Query(query, id)
	} else {
		result, err = r.db.Query(query)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrWarehouseNotFound
			return
		}
		return
	}

	var reportProducts internal.ReportProduct
	if result.Next() {
		err = result.Scan(&reportProducts.Name, &reportProducts.ProductCount)
		if err != nil {
			return
		}
		rp = append(rp, reportProducts)
	}
	return
}
