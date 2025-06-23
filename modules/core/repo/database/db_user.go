package repository

import (
	"errors"
	"fmt"

	"github.com/PeteCroton/go-basic/modules/core/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(models.UserTable) (*models.UserTable, error)
	Update(models.UserTable) (*models.UserTable, error)
	Delete(uint) error
	GetAll() ([]models.UserTable, error)
	GetById(uint) (*models.UserTable, error)
	DeletePermanently(uint) error
	FindByCredentials(loginData models.LoginRequest) (*models.UserTable, error)
}

type userRepositoryDB struct {
	db *gorm.DB
}

func NewUserRepositoryDB(db *gorm.DB) UserRepository {
	return userRepositoryDB{db: db}
}

//-------------------------------------------------------------

func (r userRepositoryDB) Create(data models.UserTable) (*models.UserTable, error) {

	data_old := models.UserTable{}
	result_old := r.db.Where("username = ?", data.Username).First(&data_old)

	if result_old.Error == nil {
		// พบผู้ใช้แล้ว = username ซ้ำ
		return nil, fmt.Errorf("username already exists")
	}

	// if result_old.Error != nil {
	// 	return nil, result_old.Error
	// }

	//Not found data...
	if errors.Is(result_old.Error, gorm.ErrRecordNotFound) {

		//เข้ารหัส Password
		if pw, err := bcrypt.GenerateFromPassword([]byte(data.Password), 0); err == nil {
			//r.db.Statement.SetColumn("Password", pw)
			data.Password = string(pw)
		}
		result := r.db.Create(&data)
		if result.Error != nil {
			return nil, result.Error
		}

		lastInsertedID := data.ID

		rs_result, err := r.GetById(lastInsertedID)
		if err != nil {
			return nil, err
		}

		return rs_result, nil
	} else {
		// error อื่นๆ (DB ล่ม, connection error ฯลฯ)
		return nil, result_old.Error
	}

}

func (r userRepositoryDB) Update(data models.UserTable) (*models.UserTable, error) {

	rs_data := models.UserTable{}
	result_search := r.db.Where("id = ?", data.ID).First(&rs_data)

	if errors.Is(result_search.Error, gorm.ErrRecordNotFound) {
		// ไม่พบข้อมูล
		return nil, fmt.Errorf("user data with User Id %d does not exist", data.ID)
	} else if result_search.Error != nil {
		// error อื่น ๆ (DB error ฯลฯ)
		return nil, result_search.Error
	}

	// เจอแล้ว → อัปเดต
	if pw, err := bcrypt.GenerateFromPassword([]byte(data.Password), 0); err == nil {
		//r.db.Statement.SetColumn("Password", pw)
		data.Password = string(pw)
	}

	result := r.db.Model(&data).Where("id = ?", data.ID).Updates(data)
	if result.Error != nil {
		return nil, result.Error
	}

	return &data, nil
}

func (r userRepositoryDB) Delete(id uint) error {
	var user models.UserTable

	// ตรวจสอบว่าข้อมูลมีอยู่หรือไม่
	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user with ID %d does not exist", id)
	} else if result.Error != nil {
		return result.Error
	}

	// ลบข้อมูล
	if err := r.db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r userRepositoryDB) GetAll() ([]models.UserTable, error) {
	listData := []models.UserTable{}

	r.db.Find(&listData)

	return listData, nil
}

func (r userRepositoryDB) GetById(id uint) (*models.UserTable, error) {

	data := models.UserTable{}
	result := r.db.Where("id = ?", id).First(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return &data, nil
}

// Example remove permanently
func (r userRepositoryDB) DeletePermanently(id uint) error {
	data := models.UserTable{}

	result := r.db.Unscoped().Where("id = ?", id).Delete(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r userRepositoryDB) FindByCredentials(loginData models.LoginRequest) (*models.UserTable, error) {
	data := models.UserTable{}
	result := r.db.Where("username = ?", loginData.Username).First(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(loginData.Password)); err != nil {
		return nil, err
	}

	return &data, nil
}
