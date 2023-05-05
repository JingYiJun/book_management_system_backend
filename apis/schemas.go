package apis

import (
	"book_management_system_backend/models"
	"time"
)

type UserInfo struct {
	IsAdmin  bool    `json:"is_admin" default:"false"`
	Avatar   *string `json:"avatar"`
	RealName *string `json:"real_name"`
	Gender   *string `json:"gender"`
	StaffID  *string `json:"staff_id"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type RegisterRequest struct {
	LoginRequest
	UserInfo
}

type TokenResponse struct {
	AccessToken string `json:"access"`
	Message     string `json:"message"`
}

type UserModifyRequest struct {
	Username *string `json:"username" validate:"omitempty,min=1"`
	Password *string `json:"password" validate:"omitempty,min=8,max=30"`
	UserInfo
}

type UserListRequest struct {
	models.PageRequest
	OrderBy string `query:"order_by" validate:"oneof=id username staff_id register_time last_login" default:"id"`
	Sort    string `query:"sort" validate:"oneof=asc desc" default:"asc"`
}

type BookListRequest struct {
	models.PageRequest
	OrderBy string  `query:"order_by" validate:"oneof=id isbn updated_at created_at title author press published_date price stock" default:"id"`
	Sort    string  `query:"sort" validate:"oneof=asc desc" default:"asc"`
	Title   *string `query:"title"`
	Author  *string `query:"author"`
	Press   *string `query:"press"`
}

func ToOrderString(orderBy string, sort string) string {
	return orderBy + " " + sort
}

type BookCreateRequest struct {
	ISBN          string     `json:"isbn" validate:"required,min=1"`
	Title         string     `json:"title" validate:"required,min=1"`
	Description   *string    `json:"description"`
	Author        string     `json:"author" validate:"required,min=1"`
	Press         string     `json:"press" validate:"required,min=1"`
	PublishedDate *time.Time `json:"published_date"`
	Price         *float64   `json:"price" validate:"omitempty,min=0"`
}

type BookModifyRequest struct {
	Title         *string    `json:"title" validate:"omitempty,min=1"`
	Description   *string    `json:"description"`
	Author        *string    `json:"author" validate:"omitempty,min=1"`
	Press         *string    `json:"press" validate:"omitempty,min=1"`
	PublishedDate *time.Time `json:"published_date"`
	Price         *float64   `json:"price" validate:"omitempty,min=0"`
	OnSale        *bool      `json:"on_sale"`
}

type PurchaseListRequest struct {
	models.PageRequest
	OrderBy string `query:"order_by" validate:"oneof=id created_at updated_at book_id user_id" default:"id"`
	Sort    string `query:"sort" validate:"oneof=asc desc" default:"asc"`
	BookID  *int   `query:"book_id"`
	UserID  *int   `query:"user_id"`
}

type PurchaseCreateRequest struct {
	BookID   int     `json:"book_id" validate:"required,min=1"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required,min=0"`
}

type PurchaseModifyRequest struct {
	Quantity *int     `json:"quantity" validate:"omitempty,min=1"`
	Price    *float64 `json:"price" validate:"omitempty,min=0"`
}
