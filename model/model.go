package model

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	//relasi
	Store Store `gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE;"`
}

type Store struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`

	//admin
	AdminID uint `gorm:"index;unique"`

	//product
	Product   []Product `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
}

type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Stock int    `gorm:"not null"`

	//store
	StoreID   uint `gorm:"index"`
	CreatedAt time.Time
}

type CartItem struct {
	ID               uint `gorm:"primaryKey"`
	PurchaseAmount   int  `gorm:"not null"`
	IsPaid           bool `gorm:"default:false"`
	CreatedAt        time.Time
	IsProductDeleted bool `gorm:"default:false"`
	//user
	UserID uint `gorm:"index"`
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	//store
	ProductID *uint    `gorm:"index"`
	Product   *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL;"`
}
