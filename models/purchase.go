package models

import "time"

type Purchase struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	BookID    int       `json:"book_id" gorm:"not null"`
	UserID    int       `json:"user_id" gorm:"not null"`
	Book      *Book     `json:"book,omitempty"`
	User      *User     `json:"user,omitempty"`
	Quantity  int       `json:"quantity" gorm:"not null;check:quantity>=1"`
	Price     float64   `json:"price" gorm:"type:numeric(10,2) not null;check:price>=0"`
	Paid      bool      `json:"paid" gorm:"default:false;not null"`
	Arrived   bool      `json:"arrived" gorm:"default:false;not null"`  // 已付款状态下可收货
	Returned  bool      `json:"returned" gorm:"default:false;not null"` // 未付款状态下可退货
}
