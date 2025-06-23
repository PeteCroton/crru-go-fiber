package models

type RoleTable struct {
	ID    uint   `json:"id"`
	Title string `json:"title" validate:"required"`
}

func (t RoleTable) TableName() string {
	return "roles"
}
