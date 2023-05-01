package models

import "time"

type Book struct {
	ID            int        `json:"id"`
	CreatedAt     time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"not null"`
	DeletedAt     time.Time  `json:"-" gorm:"index"`
	ISBN          string     `json:"isbn" gorm:"not null"`
	Title         string     `json:"title" gorm:"not null"`
	Description   *string    `json:"description"`
	Author        string     `json:"author" gorm:"not null"`
	Press         string     `json:"press"`
	PublishedDate *time.Time `json:"published_date"`
	Cover         *string    `json:"cover"` // cover url or base64, null if not set
	Price         float64    `json:"price" gorm:"type:numeric(10,2) not null"`
	Stock         int        `json:"stock" gorm:"default:0;not null"`
}

type Books []*Book
