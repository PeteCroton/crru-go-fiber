package repository

import (
	"errors"
	"fmt"

	"github.com/PeteCroton/go-basic/modules/core/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(models.RoleTable) (*models.RoleTable, error)
	Update(models.RoleTable) (*models.RoleTable, error)
	Delete(uint) error
	GetAll() ([]models.RoleTable, error)
	GetById(uint) (*models.RoleTable, error)
	DeletePermanently(uint) error
}

type roleRepositoryDB struct {
	db *gorm.DB
}

func NewRoleRepositoryDB(db *gorm.DB) RoleRepository {
	return roleRepositoryDB{db: db}
}

//-------------------------------------------------------------

func (r roleRepositoryDB) Create(data models.RoleTable) (*models.RoleTable, error) {

	result := r.db.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	lastInsertedID := data.ID

	fact_result, err := r.GetById(lastInsertedID)
	if err != nil {
		return nil, err
	}

	return fact_result, nil
}

func (r roleRepositoryDB) Update(data models.RoleTable) (*models.RoleTable, error) {
	rs_data := models.RoleTable{}
	result_search := r.db.Where("id = ?", data.ID).First(&rs_data)

	if errors.Is(result_search.Error, gorm.ErrRecordNotFound) {
		// ไม่พบข้อมูล
		return nil, fmt.Errorf("role data with ID %d does not exist", data.ID)
	} else if result_search.Error != nil {
		// error อื่น ๆ (DB error ฯลฯ)
		return nil, result_search.Error
	}

	// เจอแล้ว → อัปเดต
	result := r.db.Model(&rs_data).Updates(data)
	if result.Error != nil {
		return nil, result.Error
	}

	return &rs_data, nil
}

func (r roleRepositoryDB) Delete(id uint) error {
	// ค้นหาข้อมูลก่อนว่ามีจริงหรือไม่
	var role models.RoleTable
	result := r.db.First(&role, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("role data with ID %d does not exist", id)
	} else if result.Error != nil {
		return result.Error
	}

	// ลบข้อมูล
	if err := r.db.Delete(&role).Error; err != nil {
		return err
	}

	return nil
}

func (r roleRepositoryDB) GetAll() ([]models.RoleTable, error) {
	listData := []models.RoleTable{}

	//r.db.Find(&listData)
	r.db.Order("id ASC").Find(&listData)

	return listData, nil
}

func (r roleRepositoryDB) GetById(id uint) (*models.RoleTable, error) {

	data := models.RoleTable{}
	result := r.db.Where("id = ?", id).First(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return &data, nil
}

// Example remove permanently
func (r roleRepositoryDB) DeletePermanently(id uint) error {
	data := models.RoleTable{}

	result := r.db.Unscoped().Where("id = ?", id).Delete(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
