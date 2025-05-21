package dto

import "time"

//auth
type RegisterReq struct {
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//shop

//store
type CreateStoreReq struct {
	AdminID uint   `json:"-"`
	Name    string `json:"name"`
}

type UpdateStoreReq struct {
	UserID uint   `json:"-"`
	ID     uint   `json:"-"`
	Name   string `json:"name"`
}

type StoreAndProduct struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	AdminID   uint      `json:"admin_id"`
	CreatedAt time.Time `json:"created_at"`
	Product   []Product `json:"product"`
}

type JustStore struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	AdminID   uint      `json:"admin_id"`
	CreatedAt time.Time `json:"created_at"`
}

//product

type CreateProductReq struct {
	UserID  uint   `json:"-"`
	StoreID uint   `json:"-"`
	Name    string `json:"name"`
	Stock   int    `json:"stock"`
}

type UpdateProductReq struct {
	ID     uint   `json:"-"`
	UserID uint   `json:"-"`
	Name   string `json:"name"`
	Stock  int    `json:"stock"`
}

type Product struct {
	StoreID   uint      `json:"store_id"`
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}

//cart item
type CartItem struct {
	ID               uint `json:"id"`
	UserID           uint `json:"user_id"`
	ProductID        uint `json:"product_id"`
	PurchaseAmount   int  `json:"purchase_amount"`
	IsPaid           bool `json:"id_paid"`
	CreatedAt        time.Time
	IsProductDeleted bool `json:"is_product_deleted"`
}

type CreateCartItemReq struct {
	UserID         uint `json:"-"`
	ProductID      uint `json:"-"`
	PurchaseAmount int  `json:"purchase_amount"`
}

type UpdateAmountCartItemReq struct {
	UserID         uint `json:"-"`
	ID             uint `json:"-"`
	ProductID      uint `json:"-"`
	PurchaseAmount int  `json:"purchase_amount"`
}

type UpdatePaidCartItemReq struct {
	UserID         uint   `json:"-"`
	Email          string `json:""`
	ID             uint   `json:"-"`
	ProductID      uint   `json:"-"`
	PurchaseAmount int    `json:"purchase_amount"`
}
