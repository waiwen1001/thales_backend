package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"waiwen.com/thales-backend/models"
	"waiwen.com/thales-backend/utils"
)

type ProductController struct {
	DB *sql.DB
}

func (ac *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, "Unable to process form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// store file
	imagePath, err := utils.StoreFile(r, "image")
	if err != nil {
		http.Error(w, "Unable to store file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var product models.Product
	product.Name = r.FormValue("name")
	product.Type = r.FormValue("type")

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "Invalid price format: "+err.Error(), http.StatusBadRequest)
		return
	}

	product.Price = price
	product.Description = r.FormValue("description")
	product.PictureUrl = imagePath

	ctx := r.Context()
	if err := product.CreateProduct(ctx, ac.DB); err != nil {
		http.Error(w, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"product": product})
}

func (ac *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, "Invalid ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var product models.Product
	product.Id = id
	ctx := r.Context()
	if err := product.CheckProductExist(ctx, ac.DB); err != nil {
		// product not exist
		http.Error(w, "Product does not exist: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, "Unable to process form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// store file
	imagePath, err := utils.StoreFile(r, "image")
	if err != nil {
		http.Error(w, "Unable to store file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	product.Name = r.FormValue("name")
	product.Type = r.FormValue("type")

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "Invalid price format: "+err.Error(), http.StatusBadRequest)
		return
	}

	product.Price = price
	product.Description = r.FormValue("description")
	product.PictureUrl = imagePath

	if err := product.UpdateProduct(ctx, ac.DB); err != nil {
		http.Error(w, "Failed to update product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"product": product})
}

func (ac *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	// get pagination
	reqPaginate, err := utils.GetPaginateFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to get all products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// get sort
	query := r.URL.Query()
	sortKey := query.Get("sortKey")
	sortOrder := query.Get("sortOrder")
	search := query.Get("search")

	req := models.ProductRequest{}
	req.Page = reqPaginate.Page
	req.PageSize = reqPaginate.PageSize
	req.SortKey = sortKey
	req.SortOrder = sortOrder
	req.Search = search

	ctx := r.Context()
	product := models.Product{}
	data, paginate, err := product.GetAllProducts(req, ctx, ac.DB)
	if err != nil {
		http.Error(w, "Failed to get all products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"products": data, "paginate": paginate})
}

func (ac *ProductController) GetProductById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, "Invalid ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	product := models.Product{Id: id}
	err = product.GetProductById(ctx, ac.DB)
	if err != nil {
		http.Error(w, "Failed to get product: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"product": product})
}

func (ac *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if err != nil {
		http.Error(w, "Invalid ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	product := models.Product{Id: id}
	err = product.CheckProductExist(ctx, ac.DB)
	if err != nil {
		http.Error(w, "Product does not exist: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := product.DeleteProduct(ctx, ac.DB); err != nil {
		http.Error(w, "Failed to delete product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Product deleted successfully"})
}
