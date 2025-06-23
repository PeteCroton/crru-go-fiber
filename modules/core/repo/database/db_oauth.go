package repository

import (
	"github.com/PeteCroton/go-basic/modules/core/models"
	"gorm.io/gorm"
)

type OauthRepository interface {
	Create(models.OauthTable) (*models.OauthTable, error)
	Update(models.OauthTable) (*models.OauthTable, error)
	Delete(uint) error
	GetAll() ([]models.OauthTable, error)
	GetById(uint) (*models.OauthTable, error)
	DeletePermanently(uint) error
}

type oauthRepositoryDB struct {
	db *gorm.DB
}

func NewOauthRepositoryDB(db *gorm.DB) OauthRepository {
	return oauthRepositoryDB{db: db}
}

//-------------------------------------------------------------

func (r oauthRepositoryDB) Create(data models.OauthTable) (*models.OauthTable, error) {

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
}

func (r oauthRepositoryDB) Update(data models.OauthTable) (*models.OauthTable, error) {

	result := r.db.Model(&data).Where("id = ?", data.ID).Updates(data)
	if result.Error != nil {
		return nil, result.Error
	}

	return &data, nil
}

func (r oauthRepositoryDB) Delete(id uint) error {
	data := models.OauthTable{}

	result := r.db.Where("id = ?", id).Delete(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r oauthRepositoryDB) GetAll() ([]models.OauthTable, error) {
	listData := []models.OauthTable{}

	r.db.Find(&listData)

	return listData, nil
}

func (r oauthRepositoryDB) GetById(id uint) (*models.OauthTable, error) {

	data := models.OauthTable{}
	result := r.db.Where("id = ?", id).First(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return &data, nil
}

// Example remove permanently
func (r oauthRepositoryDB) DeletePermanently(id uint) error {
	data := models.OauthTable{}

	result := r.db.Unscoped().Where("id = ?", id).Delete(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
