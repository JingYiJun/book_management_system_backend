package models

import (
	"book_management_system_backend/utils"
	"errors"
	"gorm.io/gorm"
	"time"
)

var ErrStockNotEnough = utils.BadRequest("库存不足")
var ErrNotOnSale = utils.BadRequest("该书不在售")
var ErrBookNotFound = utils.NotFound("书籍不存在")
var ErrBookPriceNotSet = utils.BadRequest("书籍价格未设置")

type Sale struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	BookID    int       `json:"book_id" gorm:"not null"`
	UserID    int       `json:"user_id" gorm:"not null"`
	Book      *Book     `json:"-"`
	User      *User     `json:"-"`
	Quantity  int       `json:"quantity" gorm:"not null;check:quantity>=1"`
	Price     int       `json:"price" gorm:"not null;check:price>=0"` // 单价, 用 int 表示以分为单位，避免浮点数精度问题
}

func (s *Sale) PriceFloat() float64 {
	return float64(s.Price) / 100
}

func (s *Sale) BeforeCreate(tx *gorm.DB) (err error) {
	var book Book
	// Get the book
	if err = tx.Clauses(LockClause).Take(&book, s.BookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		} else {
			return
		}
	}

	// Check stock
	if book.Stock < s.Quantity {
		return ErrStockNotEnough
	}
	if !book.OnSale {
		return ErrNotOnSale
	}
	if s.Price == 0 {
		if book.Price == nil {
			return ErrBookPriceNotSet
		}
		s.Price = *book.Price
	}

	s.Book = &book
	return
}

func (s *Sale) AfterCreate(tx *gorm.DB) (err error) {
	// Update book stock
	if err = tx.Model(s.Book).Update("stock", s.Book.Stock-s.Quantity).Error; err != nil {
		return
	}
	// Create balance
	var balance = &Balance{
		UserID:        s.UserID,
		Change:        s.Price * s.Quantity,
		OperationType: OperationTypeSale,
		OperationID:   s.ID,
	}

	return tx.Create(balance).Error
}
