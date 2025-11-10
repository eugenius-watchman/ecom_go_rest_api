package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/eugenius-watchman/ecom_go_rest_api/types"
	"github.com/eugenius-watchman/ecom_go_rest_api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods("GET")      
	router.HandleFunc("/products", h.handleCreateProduct).Methods("POST")
	router.HandleFunc("/products/{id}", h.handleUpdateProduct).Methods("PUT") 
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload
	
	// Parse JSON payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	
	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %w", errors))
		return
	}
	
	// Create product in database
	err := h.store.CreateProduct(types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})
	
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Product created successfully",
	})
}

func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	// getting product from URL parameters
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	var payload types.UpdateProductPayload

	// parse payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %w", errors))
		return
	}

	// get existing product to effect or merge changes
	existingProduct, err := h.store.GetProductByID(productID) 
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	// update only fields provided 
	updatedProduct := *existingProduct
	if payload.Name != "" {
		updatedProduct.Name = payload.Name
	}
	if payload.Description != "" {
		updatedProduct.Description = payload.Description
	}
	if payload.Image != "" {
		updatedProduct.Image = payload.Image
	}
	if payload.Price > 0 {
		updatedProduct.Price = payload.Price
	}
	if payload.Quantity >= 0 {
		updatedProduct.Quantity = payload.Quantity
	}

	// update product in db
	err = h.store.UpdateProduct(productID, updatedProduct)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Product updated successfully",
	})
}