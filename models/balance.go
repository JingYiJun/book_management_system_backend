package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Balance struct {
	ID            int           `json:"id"`
	CreatedAt     time.Time     `json:"created_at"`
	Change        int           `json:"change"`                // int 表示以分为单位，避免浮点数精度问题
	Total         int           `json:"total" gorm:"not null"` // allow negative
	UserID        int           `json:"user_id" gorm:"not null"`
	User          *User         `json:"-"`
	OperationType OperationType `json:"operation_type" gorm:"not null"`
	OperationID   int           `json:"operation_id"`
	Reason        *string       `json:"reason"`
}

func (b *Balance) ChangeFloat() float64 {
	return float64(b.Change) / 100
}

func (b *Balance) TotalFloat() float64 {
	return float64(b.Total) / 100
}

func (b *Balance) Info() string {
	if b.OperationType == OperationTypeInitialize {
		return "初始化"
	}
	message := OperationTypeMap[b.OperationType]
	return fmt.Sprintf("用户 %d %s %f 元", b.UserID, message, b.ChangeFloat())
}

type OperationType = int

const (
	OperationTypePurchase OperationType = iota + 1
	OperationTypeSale
	OperationTypeManual
	OperationTypeInitialize
)

var OperationTypeMap = map[OperationType]string{
	OperationTypePurchase:   "采购支出",
	OperationTypeSale:       "销售收入",
	OperationTypeManual:     "手动收支",
	OperationTypeInitialize: "初始化",
}

func (b *Balance) BeforeCreate(tx *gorm.DB) (err error) {
	var oldBalance Balance
	// lock the last record
	if err = tx.Clauses(LockClause).Last(&oldBalance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	b.Total = oldBalance.Total + b.Change
	return nil
}
