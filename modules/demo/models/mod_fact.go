package models

import "gorm.io/gorm"

/*
gorm.Model ประกอบด้วย
- id uint
- CreatedAt time.Time
- UpdatedAt  time.Time
- DeletedAt  gorm.DeletedAt
*/

type FactTable struct {
	gorm.Model
	ID       uint   `json:"id"`
	Question string `json:"question" gorm:"text;not null;default:'-';size=100"`
	Answer   string `json:"answer" gorm:"text;not null;default:'-';size=30"`
}

// # ชื่อตารางจริง สามารถใส่ตรงนี้ หรือ ตั้งชื่อตัวแปรมาใส่ struct name ได้เลย
func (t FactTable) TableName() string {
	return "facts"
}

type FactRequest struct {
	ID       uint   `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type FactResponse struct {
	ID       uint   `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
