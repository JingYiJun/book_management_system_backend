package models

import "gorm.io/gorm"

type PageRequest struct {
	PageNum  *int `json:"page_num" query:"page_num" validate:"omitempty,min=1"`
	PageSize *int `json:"page_size" query:"page_size" validate:"omitempty,min=10,max=100"`
}

func (q PageRequest) QuerySet(tx *gorm.DB) *gorm.DB {
	if q.PageNum == nil || q.PageSize == nil {
		return tx
	}
	return tx.Offset((*(q.PageNum) - 1) * *(q.PageSize)).Limit(*(q.PageSize))
}

type Map map[string]any
