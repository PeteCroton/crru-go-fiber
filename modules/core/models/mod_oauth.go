package models

import "gorm.io/gorm"

type OauthTable struct {
	gorm.Model
	ID           uint   `json:"id" gorm:"primary_key" autoIncrement:"true"`
	UserID       uint   `json:"user_id" gorm:"not null"`
	AccessToken  string `json:"access_token" gorm:"text;not null;default:'-'"`
	RefreshToken string `json:"refresh_token" gorm:"text;not null;default:'-'"`
}


// # ชื่อตารางจริง สามารถใส่ตรงนี้ หรือ ตั้งชื่อตัวแปรมาใส่ struct name ได้เลย
func (t OauthTable) TableName() string {
	return "oauth"
}

type OauthRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type OauthResponse struct {
	ID           uint   `json:"id"`
	AccessToken  string `json:"access_token" gorm:"text;not null;default:'-'"`
	RefreshToken string `json:"refresh_token" gorm:"text;not null;default:'-'"`
	UserTable    UserTable
	UserID       uint `json:"user_id"`
}
