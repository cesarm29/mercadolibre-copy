package models

import (
	"time"
)

type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Description string    `json:"description" gorm:"size:255"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	Name         string    `json:"name" gorm:"size:255;not null"`
	Avatar       string    `json:"avatar" gorm:"size:500"`
	Phone        string    `json:"phone" gorm:"size:20"`
	RoleID       uint      `json:"role_id" gorm:"default:2"`
	Role         Role      `json:"role" gorm:"foreignKey:RoleID"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Category struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"size:255;not null"`
	Slug      string     `json:"slug" gorm:"uniqueIndex;size:255;not null"`
	Image     string     `json:"image" gorm:"size:500"`
	ParentID  *uint      `json:"parent_id"`
	Parent    *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children  []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	SellerID    uint      `json:"seller_id" gorm:"not null"`
	Seller      User      `json:"seller" gorm:"foreignKey:SellerID"`
	CategoryID  uint      `json:"category_id" gorm:"not null"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"not null;default:0"`
	Images      []string  `json:"images" gorm:"type:text[]"`
	Condition   string    `json:"condition" gorm:"size:20;not null;default:new"` // new, used
	Status      string    `json:"status" gorm:"size:20;not null;default:active"` // active, sold, inactive
	SoldCount   int       `json:"sold_count" gorm:"default:0"`
	FreeShipping bool     `json:"free_shipping" gorm:"default:false"`
	Location    string    `json:"location" gorm:"size:255"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CartItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Order struct {
	ID              uint        `json:"id" gorm:"primaryKey"`
	UserID          uint        `json:"user_id" gorm:"not null"`
	User            User        `json:"user" gorm:"foreignKey:UserID"`
	Total           float64     `json:"total" gorm:"not null"`
	Status          string      `json:"status" gorm:"size:20;not null;default:pending"` // pending, paid, shipped, delivered, cancelled
	ShippingAddress string      `json:"shipping_address" gorm:"type:text;not null"`
	Items           []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
}

type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	Rating    int       `json:"rating" gorm:"not null"`
	Comment   string    `json:"comment" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required,min=2"`
	Phone    string `json:"phone"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type PaginationQuery struct {
	Page    int    `form:"page,default=1"`
	Limit   int    `form:"limit,default=20"`
	Search  string `form:"search"`
	Sort    string `form:"sort,default=created_at"`
	Order   string `form:"order,default=desc"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
