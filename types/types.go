package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type ProductStore interface {
	GetProducts() ([]Product, error)
	GetProductByID(id int) (*Product, error)
	CreateProduct(Product) error//
	UpdateProduct(id int, product Product) error
	ProductExists(id int) (bool, error)
	UpdateProductQuantity(id int, newQuantity int) error
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

type User struct {
	// Go field name ... JSON field nam
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterUserPayload struct {
	// Go field name ... JSON field nam
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	// Go field name ... JSON field name
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Image       string  `json:"image" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}

type UpdateProductPayload struct {
	Name        string  `json:"name" validate:"omitempty"`
	Description string  `json:"description" validate:"omitempty"`
	Image       string  `json:"image" validate:"omitempty,url"`
	Price       float64 `json:"price" validate:"omitempty,min=0"`
	Quantity    int     `json:"quantity" validate:"omitempty,min=0"`
}

type CartStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
	GetOrderByID(id int) (*Order, error)
	GetOrdersByUserID(userID int) ([]Order, error)
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"` // paid pending shipped
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"orderId"`
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // price at time of purchase
}

// checkout
type CheckoutPayload struct {
	Address string         `json:"address" validate:"required"`
	Items   []CheckoutItem `json:"items" validate:"required,min=1"`
}

type CheckoutItem struct {
	ProductID int `json:"productId" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}
