package repository

import (
	"app/internal"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type RepositoryProductMySql struct {
	// db is the underlying database.
	db *sql.DB
}

func NewRepositoryProductMySql(db *sql.DB) (r *RepositoryProductMySql) {
	r = &RepositoryProductMySql{
		db: db,
	}
	return
}

func (r *RepositoryProductMySql) FindById(id int) (p internal.Product, err error) {
	query := "SELECT `id`, `name`, `quantity`, `code_value`, `is_published`, `price`,`expiration` FROM `products` WHERE `id` = ?"

	result := r.db.QueryRow(query, id)
	if result.Err() != nil {
		if errors.Is(result.Err(), sql.ErrNoRows) {
			err = internal.ErrRepositoryProductNotFound
			return
		}
		err = result.Err()
		return
	}

	var timeString string
	err = result.Scan(&p.Id, &p.Name, &p.Quantity, &p.CodeValue, &p.IsPublished, &p.Price, &timeString)
	if err != nil {
		return
	}

	p.Expiration, err = time.Parse(time.DateOnly, timeString)
	if err != nil {
		return
	}

	return
}

func (r *RepositoryProductMySql) Save(p *internal.Product) (err error) {
	query := "INSERT INTO `products` (`name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`, `id_warehouse`) VALUES (?, ?, ?, ?, ?, ?,1)"

	result, err := r.db.Exec(query, p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price)
	if err != nil {
		return
	}

	lastId, err := result.LastInsertId()

	(*p).Id = int(lastId)
	return
}

func (r *RepositoryProductMySql) UpdateOrSave(p *internal.Product) (err error) {
	return
}

func (r *RepositoryProductMySql) Update(p *internal.Product) (err error) {
	query := "UPDATE `products` SET `name` = ?, `quantity` = ?, `code_value` = ?, `is_published` = ?, `expiration` = ?, `price` = ? WHERE `id` = ?"

	result, err := r.db.Exec(query, p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price, p.Id)
	if err != nil {
		var mySqlErr *mysql.MySQLError
		if errors.As(err, &mySqlErr) {
			switch mySqlErr.Number {
			case 1062:
				err = fmt.Errorf("product already exists")
				return
			default:
				err = fmt.Errorf("unknown error")
				return
			}
		}
		return
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowAffected == 0 {
		err = internal.ErrRepositoryProductNotFound
		return
	}

	return
}

func (r *RepositoryProductMySql) Delete(id int) (err error) {
	query := "DELETE FROM `products` WHERE `id` = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowAffected == 0 {
		err = internal.ErrRepositoryProductNotFound
		return
	}
	return
}
