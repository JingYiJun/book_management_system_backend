package models

import "gorm.io/gorm"

type PageRequest struct {
	PageNum  int `json:"page_num" query:"page_num" validate:"required,min=1"`
	PageSize int `json:"page_size" query:"page_size" validate:"required,min=10,max=100"`
}

func (q PageRequest) QuerySet(tx *gorm.DB) *gorm.DB {
	return tx.Offset((q.PageNum - 1) * q.PageSize).Limit(q.PageSize)
}

type Map map[string]any
