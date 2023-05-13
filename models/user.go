package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type User struct {
	ID             int            `json:"id"`
	RegisterTime   time.Time      `json:"register_time" gorm:"autoCreateTime;not null"`
	LastLogin      time.Time      `json:"last_login" gorm:"autoUpdateTime;not null"`
	DeletedAt      gorm.DeletedAt `json:"-"`
	Username       string         `json:"username" gorm:"size:256;uniqueIndex;not null"`
	HashedPassword string         `json:"-" gorm:"size:256;not null"`
	IsAdmin        bool           `json:"is_admin" gorm:"default:false;not null"`
	Avatar         *string        `json:"avatar"`
	RealName       *string        `json:"real_name"`
	Gender         *string        `json:"gender" gorm:"size:1"`
	StaffID        *string        `json:"staff_id" gorm:"size:20"`
}

type UserJwtSecret struct {
	ID     int    `json:"id"`
	Secret string `json:"secret" gorm:"size:256"`
}

func (user *User) BeforeDelete(tx *gorm.DB) error {
	return tx.Model(&user).Update("username", user.Username+"_d_"+strconv.FormatInt(time.Now().Unix(), 10)).Error
}
