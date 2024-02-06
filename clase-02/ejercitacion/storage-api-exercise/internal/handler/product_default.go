package handler

import (
	"app/internal"
	"app/platform/web/request"
	"app/platform/web/response"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// NewProductsDefault returns a new instance of ProductsDefault
func NewProductsDefault(rp internal.RepositoryProducts) *ProductsDefault {
	return &ProductsDefault{
		rp: rp,
	}
}

// ProductsDefault is a struct that represents the default product handler
type ProductsDefault struct {
	// rp is the product repository
	rp internal.RepositoryProducts
}

// ProductJSON is a struct that represents a product in JSON
type ProductJSON struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
	WarehouseId int     `json:"warehouse_id"`
}

func (h *ProductsDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := h.rp.GetAll()
		if err != nil {
			if errors.Is(err, internal.ErrProductNotFound) {
				response.Error(w, http.StatusNotFound, "products not found")
				return
			}
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		var productsJSON []ProductJSON
		for _, p := range products {
			productsJSON = append(productsJSON, ProductJSON{
				ID:          p.ID,
				Name:        p.Name,
				Quantity:    p.Quantity,
				CodeValue:   p.CodeValue,
				IsPublished: p.IsPublished,
				Expiration:  p.Expiration.Format(time.DateOnly),
				Price:       p.Price,
				WarehouseId: p.WarehouseId,
			})
		}

		response.JSON(w, http.StatusOK, map[string]any{"message": "products found", "data": productsJSON})
	}
}

// GetOne returns a product by id
func (h *ProductsDefault) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		p, err := h.rp.GetOne(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Error(w, http.StatusNotFound, "product not found")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error"+err.Error())
			}
			return
		}

		// response
		// - serialize
		data := ProductJSON{
			ID:          p.ID,
			Name:        p.Name,
			Quantity:    p.Quantity,
			CodeValue:   p.CodeValue,
			IsPublished: p.IsPublished,
			Expiration:  p.Expiration.Format(time.DateOnly),
			Price:       p.Price,
			WarehouseId: p.WarehouseId,
		}
		response.JSON(w, http.StatusOK, map[string]any{"message": "product found", "data": data})
	}
}

// RequestBodyProductCreate is a struct that represents the request body of a product to create
type RequestBodyProductCreate struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
	WarehouseId int     `json:"warehouse_id"`
}

// Create creates a product
func (h *ProductsDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var body RequestBodyProductCreate
		if err := request.JSON(r, &body); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}
		exp, err := time.Parse(time.DateOnly, body.Expiration)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid expiration date")
			return
		}

		// process
		p := internal.Product{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  exp,
			Price:       body.Price,
			WarehouseId: body.WarehouseId,
		}
		if err := h.rp.Store(&p); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotUnique):
				response.Error(w, http.StatusConflict, "product not unique")
			case errors.Is(err, internal.ErrProductRelation):
				response.Error(w, http.StatusConflict, "product relation error")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// - serialize
		data := ProductJSON{
			ID:          p.ID,
			Name:        p.Name,
			Quantity:    p.Quantity,
			CodeValue:   p.CodeValue,
			IsPublished: p.IsPublished,
			Expiration:  p.Expiration.Format(time.DateOnly),
			Price:       p.Price,
			WarehouseId: p.WarehouseId,
		}
		response.JSON(w, http.StatusCreated, map[string]any{"message": "product created", "data": data})
	}
}

// RequestBodyProductUpdate is a struct that represents the request body of a product to update
type RequestBodyProductUpdate struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
	WarehouseId int     `json:"warehouse_id"`
}

// Update updates a product
func (h *ProductsDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		// - get product
		p, err := h.rp.GetOne(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Error(w, http.StatusNotFound, "product not found")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}
		// - patch product
		body := RequestBodyProductUpdate{
			Name:        p.Name,
			Quantity:    p.Quantity,
			CodeValue:   p.CodeValue,
			IsPublished: p.IsPublished,
			Expiration:  p.Expiration.Format(time.DateOnly),
			Price:       p.Price,
			WarehouseId: p.WarehouseId,
		}
		if err := request.JSON(r, &body); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}
		exp, err := time.Parse(time.DateOnly, body.Expiration)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid expiration date")
			return
		}
		p.Name = body.Name
		p.Quantity = body.Quantity
		p.CodeValue = body.CodeValue
		p.IsPublished = body.IsPublished
		p.Expiration = exp
		p.Price = body.Price
		p.WarehouseId = body.WarehouseId
		// - update product
		if err := h.rp.Update(&p); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotUnique):
				response.Error(w, http.StatusConflict, "product not unique")
			case errors.Is(err, internal.ErrProductRelation):
				response.Error(w, http.StatusConflict, "product relation error")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// - serialize
		data := ProductJSON{
			ID:          p.ID,
			Name:        p.Name,
			Quantity:    p.Quantity,
			CodeValue:   p.CodeValue,
			IsPublished: p.IsPublished,
			Expiration:  p.Expiration.Format(time.DateOnly),
			Price:       p.Price,
			WarehouseId: p.WarehouseId,
		}
		response.JSON(w, http.StatusOK, map[string]any{"message": "product updated", "data": data})
	}
}

// Delete deletes a product by id
func (h *ProductsDefault) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		if err := h.rp.Delete(id); err != nil {
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// response
		response.JSON(w, http.StatusOK, map[string]any{"message": "product deleted", "data": id})
	}
}
