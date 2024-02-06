package handler

import (
	"app/internal"
	"app/platform/web/response"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WarehouseDefault struct {
	rp internal.WarehouseRepository
}

func NewWarehouseDefault(rp internal.WarehouseRepository) *WarehouseDefault {
	return &WarehouseDefault{
		rp: rp,
	}
}

type WarehouseJSON struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Capacity  int    `json:"capacity"`
}

type BodyWarehouseJSON struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Capacity  int    `json:"capacity"`
}

type ReportProduct struct {
	Name         string `json:"name"`
	ProductCount int    `json: "product_count"`
}

func (h *WarehouseDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouses, err := h.rp.GetAll()
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseNotFound):
				response.Error(w, http.StatusNotFound, "warehouse not found")
				return
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error"+err.Error())
				return
			}
		}

		var wawrehousesJSON []WarehouseJSON
		for _, w := range warehouses {
			wawrehousesJSON = append(wawrehousesJSON, WarehouseJSON{
				Id:        w.Id,
				Name:      w.Name,
				Address:   w.Address,
				Telephone: w.Telephone,
				Capacity:  w.Capacity,
			})
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message":    "warehouses found",
			"warehouses": wawrehousesJSON,
		})
	}
}
func (h *WarehouseDefault) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// requests
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		warehouse, err := h.rp.GetOne(idInt)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseNotFound):
				response.Error(w, http.StatusNotFound, "warehouse not found")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// serialize response
		warehouseJSON := WarehouseJSON{
			Id:        warehouse.Id,
			Name:      warehouse.Name,
			Address:   warehouse.Address,
			Telephone: warehouse.Telephone,
			Capacity:  warehouse.Capacity,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message":   "warehouse found",
			"warehouse": warehouseJSON,
		})
	}
}

func (h *WarehouseDefault) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// requests
		var warehouseJSON BodyWarehouseJSON
		err := json.NewDecoder(r.Body).Decode(&warehouseJSON)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request")
			return
		}

		//serialize request
		warehouse := internal.Warehouse{
			Name:      warehouseJSON.Name,
			Address:   warehouseJSON.Address,
			Telephone: warehouseJSON.Telephone,
			Capacity:  warehouseJSON.Capacity,
		}
		err = h.rp.Store(&warehouse)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseAlreadyExists):
				response.Error(w, http.StatusConflict, "warehouse already exists")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// serialize response
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "warehouse created",
			"data":    warehouseJSON,
		})

	}
}

func (h *WarehouseDefault) ReportProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			idInt = 0
		}

		rp, err := h.rp.ReportProducts(idInt)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseNotFound):
				response.Error(w, http.StatusNotFound, "warehouse not found")
				return
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error"+err.Error())
				return
			}
		}

		var reportsProduct []ReportProduct
		for _, v := range rp {
			reportsProduct = append(reportsProduct, ReportProduct{
				Name:         v.Name,
				ProductCount: v.ProductCount,
			})
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "generate report product success",
			"data":    reportsProduct,
		})

	}
}
