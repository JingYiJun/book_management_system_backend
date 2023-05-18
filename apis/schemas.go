package apis

import (
	"book_management_system_backend/models"
	"time"
)

func ToOrderString(orderBy string, sort string) string {
	return orderBy + " " + sort
}

type UserInfo struct {
	IsAdmin  bool    `json:"is_admin" default:"false"`
	Avatar   *string `json:"avatar,omitempty"`
	RealName *string `json:"real_name,omitempty"`
	Gender   *string `json:"gender,omitempty"`
	StaffID  *string `json:"staff_id,omitempty"`
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
	OrderBy string `json:"order_by" query:"order_by" validate:"oneof=id username staff_id register_time last_login" default:"id"`
	Sort    string `json:"sort" query:"sort" validate:"oneof=asc desc" default:"asc"`
}

type UserResponse struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	RegisterTime time.Time `json:"register_time"`
	LastLogin    time.Time `json:"last_login"`
	UserInfo
}

type UserListResponse struct {
	Users     []UserResponse `json:"users"`
	PageTotal int            `json:"page_total"`
}

/* Book */

type BookListRequest struct {
	models.PageRequest
	OrderBy string  `json:"order_by" query:"order_by" validate:"oneof=id isbn updated_at created_at title author press published_date price stock" default:"id"`
	Sort    string  `json:"sort" query:"sort" validate:"oneof=asc desc" default:"asc"`
	Title   *string `json:"title" query:"title"`
	Author  *string `json:"author" query:"author"`
	Press   *string `json:"press" query:"press"`
	OnSale  *bool   `json:"on_sale" query:"on_sale"`
	ID      *int    `json:"id" query:"id"`
	ISBN    *string `json:"isbn" query:"isbn"`
}

type BookCreateRequest struct {
	ISBN          string     `json:"isbn" validate:"required,min=1"`
	Title         string     `json:"title" validate:"required,min=1"`
	Description   *string    `json:"description"`
	Author        string     `json:"author" validate:"required,min=1"`
	Press         string     `json:"press" validate:"required,min=1"`
	PublishedDate *time.Time `json:"published_date"`
	PriceFloat    *float64   `json:"price" validate:"omitempty,min=0"`
	Cover         *string    `json:"cover"` // cover url or base64, null if not set
	OnSale        bool       `json:"on_sale" default:"false"`
}

func (b *BookCreateRequest) Price() *int {
	if b.PriceFloat == nil {
		return nil
	}
	price := int(*b.PriceFloat * 100)
	return &price
}

type BookModifyRequest struct {
	Title         *string    `json:"title" validate:"omitempty,min=1"`
	Description   *string    `json:"description"`
	Author        *string    `json:"author" validate:"omitempty,min=1"`
	Press         *string    `json:"press" validate:"omitempty,min=1"`
	PublishedDate *time.Time `json:"published_date"`
	PriceFloat    *float64   `json:"price" validate:"omitempty,min=0"`
	Cover         *string    `json:"cover"` // cover url or base64, null if not set
	OnSale        *bool      `json:"on_sale"`
}

func (b *BookModifyRequest) Price() *int {
	if b.PriceFloat == nil {
		return nil
	}
	price := int(*b.PriceFloat * 100)
	return &price
}

type BookResponse struct {
	ID            int        `json:"id"`
	CreatedAt     time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"not null"`
	UserID        int        `json:"user_id" gorm:"not null"` // user who create the book
	ISBN          string     `json:"isbn" gorm:"not null"`
	Title         string     `json:"title" gorm:"not null"`
	Description   *string    `json:"description"`
	Author        string     `json:"author" gorm:"not null"`
	Press         string     `json:"press" gorm:"not null"`
	PublishedDate *time.Time `json:"published_date"`
	Cover         *string    `json:"cover"` // cover url or base64, null if not set
	PriceFloat    *float64   `json:"price"` // 单价, 用 int 表示以分为单位，避免浮点数精度问题
	Stock         int        `json:"stock" gorm:"default:0;not null"`
	OnSale        bool       `json:"on_sale" gorm:"default:false;not null"`
}

type BookListResponse struct {
	Books     []BookResponse `json:"books"`
	PageTotal int            `json:"page_total"`
}

/* Purchase */

type PurchaseListRequest struct {
	models.PageRequest
	OrderBy string `json:"order_by" query:"order_by" validate:"oneof=id created_at updated_at book_id user_id" default:"id"`
	Sort    string `json:"sort" query:"sort" validate:"oneof=asc desc" default:"asc"`
	BookID  *int   `json:"book_id" query:"book_id"`
	UserID  *int   `json:"user_id" query:"user_id"`
}

type PurchaseCreateRequest struct {
	BookID     int     `json:"book_id" validate:"required,min=1"`
	Quantity   int     `json:"quantity" validate:"required,min=1"`
	PriceFloat float64 `json:"price" validate:"required,min=0"`
}

func (p *PurchaseCreateRequest) Price() int {
	return int(p.PriceFloat * 100)
}

type PurchaseModifyRequest struct {
	Quantity   *int     `json:"quantity" validate:"omitempty,min=1"`
	PriceFloat *float64 `json:"price" validate:"omitempty,min=0"`
}

func (p *PurchaseModifyRequest) Price() *int {
	if p.PriceFloat == nil {
		return nil
	}
	price := int(*p.PriceFloat * 100)
	return &price
}

type PurchaseResponse struct {
	ID         int           `json:"id"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	BookID     int           `json:"book_id"`
	UserID     int           `json:"user_id"`
	Quantity   int           `json:"quantity"`
	PriceFloat float64       `json:"price"`
	Paid       bool          `json:"paid"`
	Arrived    bool          `json:"arrived"`
	Returned   bool          `json:"returned"`
	Book       *BookResponse `json:"book,omitempty"`
}

type PurchaseListResponse struct {
	Purchases []PurchaseResponse `json:"purchases"`
	PageTotal int                `json:"page_total"`
}

/* Balance */

type BalanceListRequest struct {
	models.PageRequest
	OrderBy   string     `json:"order_by" query:"order_by" validate:"oneof=id created_at user_id change" default:"id"`
	Sort      string     `json:"sort" query:"sort" validate:"oneof=asc desc" default:"asc"`
	UserID    *int       `json:"user_id" query:"user_id"`
	Positive  *bool      `json:"positive" query:"positive"` // true: positive, false: negative, nil: all
	StartTime *time.Time `json:"start_time" query:"start_time"`
	EndTime   *time.Time `json:"end_time" query:"end_time"`
}

type BalanceCreateRequest struct {
	ChangeFloat float64 `json:"change" validate:"required"`
	Reason      *string `json:"reason"`
}

func (b *BalanceCreateRequest) Change() int {
	return int(b.ChangeFloat * 100)
}

type BalanceResponse struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UserID        int       `json:"user_id"`
	Change        float64   `json:"change" copier:"ChangeFloat"`
	Total         float64   `json:"balance" copier:"TotalFloat"`
	OperationType int       `json:"operation_type"`
	OperationID   int       `json:"operation_id"`
	Info          string    `json:"info"`
}

type BalanceListResponse struct {
	Balances  []BalanceResponse `json:"balances"`
	PageTotal int               `json:"page_total"`
}

/* Sale */

type SaleListRequest struct {
	models.PageRequest
	OrderBy   string     `json:"order_by" query:"order_by" validate:"oneof=id created_at updated_at book_id user_id" default:"id"`
	Sort      string     `json:"sort" query:"sort" validate:"oneof=asc desc" default:"asc"`
	BookID    *int       `json:"book_id" query:"book_id"`
	UserID    *int       `json:"user_id" query:"user_id"`
	StartTime *time.Time `json:"start_time" query:"start_time"`
	EndTime   *time.Time `json:"end_time" query:"end_time"`
}

type SaleCreateRequest struct {
	BookID     int     `json:"book_id" validate:"required,min=1"`
	Quantity   int     `json:"quantity" validate:"required,min=1"`
	PriceFloat float64 `json:"price"`
}

func (s *SaleCreateRequest) Price() int {
	return int(s.PriceFloat * 100)
}

type SaleResponse struct {
	ID         int       `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	BookID     int       `json:"book_id"`
	UserID     int       `json:"user_id"`
	Quantity   int       `json:"quantity"`
	PriceFloat float64   `json:"price"`
}

type SaleListResponse struct {
	Sales     []SaleResponse `json:"sales"`
	PageTotal int            `json:"page_total"`
}
