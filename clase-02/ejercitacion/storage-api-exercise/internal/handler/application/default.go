package application

import (
	"app/internal/handler"
	"app/internal/repository"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
)

// ConfigDefault is a struct that represents the default application configuration
type ConfigDefault struct {
	// Database is the database configuration
	Database mysql.Config
	// Address is the address of the application
	Address string
}

// NewDefault returns a new default application
func NewDefault(cfg *ConfigDefault) *Default {
	// default
	cfgDefault := &ConfigDefault{
		Address: ":8080",
	}
	if cfg != nil {
		cfgDefault.Database = cfg.Database
		if cfg.Address != "" {
			cfgDefault.Address = cfg.Address
		}
	}

	return &Default{
		cfgDb: cfgDefault.Database,
		addr:  cfgDefault.Address,
	}
}

// Default is a struct that represents the default application
type Default struct {
	// cfgDb is the database configuration
	cfgDb mysql.Config
	// addr is the address of the application
	addr string
}

// Run runs the default application
func (d *Default) Run() (err error) {
	// dependencies
	// - database: connection
	db, err := sql.Open("mysql", d.cfgDb.FormatDSN())
	if err != nil {
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return
	}
	// - router: chi
	rt := chi.NewRouter()
	// - router: middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)

	// routes
	// - product
	routesProduct(rt, db)
	// - warehouses
	routesWarehouse(rt, db)
	// run
	err = http.ListenAndServe(d.addr, rt)
	if err != nil {
		return
	}
	return
}

func routesProduct(rt *chi.Mux, db *sql.DB) {
	// - repository: products
	rp := repository.NewProductsMySQL(db)

	// - handler: products
	hp := handler.NewProductsDefault(rp)

	// - router: routes
	rt.Route("/products", func(r chi.Router) {

		// - GET all products
		r.Get("/", hp.GetAll())

		// - GET product by id
		r.Get("/{id}", hp.GetOne())

		// - POST /products
		r.Post("/", hp.Create())

		// - PUT /products/{id}
		r.Patch("/{id}", hp.Update())

		// - DELETE /products/{id}
		r.Delete("/{id}", hp.Delete())
	})
}

func routesWarehouse(rt *chi.Mux, db *sql.DB) {
	// - repository: warehouses
	rp := repository.NewWarehouseMySQL(db)

	// - handler: warehouses
	hp := handler.NewWarehouseDefault(rp)

	rt.Route("/warehouses", func(r chi.Router) {
		// - GET /warehouses
		r.Get("/", hp.GetAll())
		r.Get("/{id}", hp.GetOne())
		r.Get("/reportProducts", hp.ReportProduct())
		// - POST /warehouses
		r.Post("/", hp.Store())

	})
}
