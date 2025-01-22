package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"waiwen.com/thales-backend/controllers"
)

func InitRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// file server
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Controllers
	productController := &controllers.ProductController{DB: db}

	// API Routes
	api := router.PathPrefix("/api").Subrouter()
	product := api.PathPrefix("/products").Subrouter()

	// Routes
	product.HandleFunc("", productController.CreateProduct).Methods("POST")
	product.HandleFunc("/{id:[0-9]+}", productController.UpdateProduct).Methods("PUT")
	product.HandleFunc("", productController.GetAllProducts).Methods("GET")
	product.HandleFunc("/{id:[0-9]+}", productController.GetProductById).Methods("GET")
	product.HandleFunc("/{id:[0-9]+}", productController.DeleteProduct).Methods("DELETE")

	return router
}
