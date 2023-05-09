package models

import "time"

type Book struct {
	ID            int        `json:"id"`
	CreatedAt     time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"not null"`
	DeletedAt     time.Time  `json:"-" gorm:"index"`
	UserID        int        `json:"user_id" gorm:"not null"` // user who create the book
	User          *User      `json:"-"`
	ISBN          string     `json:"isbn" gorm:"not null"`
	Title         string     `json:"title" gorm:"not null"`
	Description   *string    `json:"description"`
	Author        string     `json:"author" gorm:"not null"`
	Press         string     `json:"press" gorm:"not null"`
	PublishedDate *time.Time `json:"published_date"`
	Cover         *string    `json:"cover"` // cover url or base64, null if not set
	Price         *int       `json:"price"` // 单价, 用 int 表示以分为单位，避免浮点数精度问题
	Stock         int        `json:"stock" gorm:"default:0;not null"`
	OnSale        bool       `json:"on_sale" gorm:"default:false;not null"`
}

func (b *Book) PriceFloat() float64 {
	if b.Price == nil {
		return 0
	}
	return float64(*b.Price) / 100
}
