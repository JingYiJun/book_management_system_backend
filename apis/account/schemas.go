package account

import "book_management_system_backend/models"

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
	OrderBy string `json:"order_by" query:"order_by" validate:"oneof=id username staff_id register_time last_login" default:"id"`
	Sort    string `json:"sort" query:"sort" validate:"oneof=asc desc" default:"asc"`
}
