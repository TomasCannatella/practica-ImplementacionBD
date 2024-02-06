package handler_test

import (
	"app/internal/handler"
	"app/internal/repository"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func init() {
	cfg := mysql.Config{
		User:                 "user1",
		Passwd:               "secret_password",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "my_db_test",
		ParseTime:            true,
		AllowNativePasswords: true, // Enable MySQL native password authentication
	}
	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func TestWarehouseDefault_GetAll(t *testing.T) {
	t.Run("success 01 - warehouse found", func(t *testing.T) {
		// arrange
		db, err := sql.Open("txdb", os.Getenv("DB_NAME_TEST"))
		require.NoError(t, err)

		err = func() error {
			_, err := db.Exec("INSERT INTO `warehouses` (`id`, `name`, `adress`, `telephone`, `capacity`) VALUES (1, 'warehouse 1', 'address 1', 'telephone 1', 100)")
			return err
		}()
		require.NoError(t, err)

		rp := repository.NewWarehouseMySQL(db)
		hd := handler.NewWarehouseDefault(rp)

		//act
		req := httptest.NewRequest("GET", "/warehouses", nil)
		res := httptest.NewRecorder()
		hd.GetAll()(res, req)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"warehouses found", "warehouses":[{"address":"address 1", "capacity":100, "id":1, "name":"warehouse 1", "telephone":"telephone 1"}]}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})
}

func TestWarehouseDefault_GetByID(t *testing.T) {
	t.Run("success 01 - warehouse found", func(t *testing.T) {
		// arrange
		db, err := sql.Open("txdb", os.Getenv("DB_NAME_TEST"))
		require.NoError(t, err)

		err = func() error {
			_, err := db.Exec("INSERT INTO `warehouses` (`id`, `name`, `adress`, `telephone`, `capacity`) VALUES (1, 'warehouse 1', 'address 1', 'telephone 1', 100)")
			return err
		}()
		require.NoError(t, err)

		rp := repository.NewWarehouseMySQL(db)
		hd := handler.NewWarehouseDefault(rp)

		req := httptest.NewRequest("GET", "/warehouses/1", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()
		hd.GetOne()(res, req)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"warehouse found", "warehouse":{"address":"address 1", "capacity":100, "id":1, "name":"warehouse 1", "telephone":"telephone 1"}}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})
}

func TestWarehouseDefault_Store(t *testing.T) {
	t.Run("success 01 - warehouse stored", func(t *testing.T) {
		// arrange
		db, err := sql.Open("txdb", os.Getenv("DB_NAME_TEST"))
		require.NoError(t, err)

		rp := repository.NewWarehouseMySQL(db)
		hd := handler.NewWarehouseDefault(rp)

		req := httptest.NewRequest("POST", "/warehouses", strings.NewReader(`{"name":"warehouse 1", "address":"address 1", "telephone":"telephone 1", "capacity":100}`))
		res := httptest.NewRecorder()
		hd.Store()(res, req)

		// assert
		expectedCode := http.StatusCreated
		expectedBody := `{
			"data": {
				"name": "warehouse 1",
				"address": "address 1",
				"telephone": "telephone 1",
				"capacity": 100
			},
			"message": "warehouse created"
		}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})
}
