package cart

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eugenius-watchman/ecom_go_rest_api/types"
	"github.com/eugenius-watchman/ecom_go_rest_api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.CartStore
	productStore types.ProductStore
}

func NewHandler(store types.CartStore, productStore types.ProductStore) *Handler {
	return &Handler{
		store:        store,
		productStore: productStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", h.handleCheckout).Methods("POST")
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	// get user id from JWT token
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))

		return
	}

	var payload types.CheckoutPayload

	// parse JSON payload
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

	// calculate total ... will need product prices from db
	total, productPrices, err := h.calculateTotalWithPrices(payload.Items)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)

		return
	}

	// create order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:    userID,
		Total:     total,
		Status:    "pending",
		Address:   payload.Address,
		CreatedAt: time.Now(),
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)

		return
	}

	// create order items
	for _, item := range payload.Items {
		// get price from productPrices map
		price := productPrices[item.ProductID]

		err := h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price, // from products db
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)

			return
		}
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Order created successfully",
		"orderId": orderID,
		"total":   total,
		"status":  "pending",
	})
}


// helper func to get actual prices from db
func (h *Handler) calculateTotalWithPrices(items []types.CheckoutItem) (float64, map[int]float64, error) {
	var total float64
	productPrices := make(map[int]float64)

	// get actual product price(s) from db
	for _, item := range items {
		product, err := h.productStore.GetProductByID(item.ProductID)

		if err != nil {
			return 0, nil, fmt.Errorf("product with ID %d not found", item.ProductID)
		}

		// check for sufficient quatity of product
		if product.Quantity < item.Quantity {
			return 0, nil, fmt.Errorf("insufficient quantity for priduct %s", product.Name)
		}

		itemTotal := product.Price * float64(item.Quantity)
		total += itemTotal
		productPrices[item.ProductID] = product.Price // store price for order items
	}

	return total, productPrices, nil
}

// helper function to get uer ID from JWT
func (h *Handler) getUserIDFromToken(r *http.Request) (int, error) {
	// implement JWT token extraction ...depends on auth middleware
	// tentatively return a dummy user ID
	return 1, nil
}

// helper function to calculate order total
func (h *Handler) calculateTotal(items []types.CheckoutItem) (float64, error) {
	var total float64

	// query product prices from db
	// using place holder prices tentatively
	for _, item := range items {
		// in a real implementation, we get price from product table
		productPrice := 10.0 // place holder
		total += productPrice * float64(item.Quantity)
	}

	return total, nil
}
