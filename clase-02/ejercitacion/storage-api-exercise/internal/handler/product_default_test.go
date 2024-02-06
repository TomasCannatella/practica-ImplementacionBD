package handler_test

import (
	"app/internal/handler"
	"app/internal/repository"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProductDefault_GetAll(t *testing.T) {
	t.Run("success 01 - products found", func(t *testing.T) {
		// arrange
		db, err := sql.Open("txdb", "my_db_test")
		require.NoError(t, err)

		err = func() error {
			_, err := db.Exec("INSERT INTO `products` (`id`, `name`, `quantity`,`code_value`,`is_published`,`expiration`,`price`,`id_warehouse`) VALUES (1, 'product 1', 100, 'code_value 1', true, '2021-12-31', 100, 1)")
			return err
		}()
		require.NoError(t, err)

		rp := repository.NewProductsMySQL(db)
		hd := handler.NewProductsDefault(rp)

		//act
		req := httptest.NewRequest("GET", "/products", nil)
		res := httptest.NewRecorder()
		hd.GetAll()(res, req)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{ "data": [{"id": 1,"name": "product 1","quantity": 100,"code_value": "code_value 1","is_published": true,"expiration": "2021-12-31","price": 100,"warehouse_id": 1}],"message": "products found"}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})

}
