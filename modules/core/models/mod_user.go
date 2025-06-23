package models

import "gorm.io/gorm"

type UserTable struct {
	gorm.Model
	Username string `json:"username" validate:"required" `
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"email,required"`
	//RolesTable RoleTable
	Role_ID uint `json:"role_id" validate:"number"`
}

// # ชื่อตารางจริง สามารถใส่ตรงนี้ หรือ ตั้งชื่อตัวแปรมาใส่ struct name ได้เลย
func (t UserTable) TableName() string {
	return "users"
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
