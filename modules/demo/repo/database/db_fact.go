package repository

import (
	"github.com/PeteCroton/go-basic/modules/demo/models"
	"gorm.io/gorm"
)

type FactRepository interface {
	Create(models.FactTable) (*models.FactTable, error)
	Update(models.FactTable) (*models.FactTable, error)
	Delete(uint) error
	GetAll() ([]models.FactTable, error)
	GetById(uint) (*models.FactTable, error)
	DeletePermanently(uint) error
}

type factRepositoryDB struct {
	db *gorm.DB //configs.DB.Db
}

func NewFactRepositoryDB(db *gorm.DB) FactRepository {
	return factRepositoryDB{db: db}
}

//-------------------------------------------------------------

func (r factRepositoryDB) Create(fact models.FactTable) (*models.FactTable, error) {

	result := r.db.Create(&fact)
	if result.Error != nil {
		return nil, result.Error
	}

	lastInsertedID := fact.ID

	fact_result, err := r.GetById(lastInsertedID)
	if err != nil {
		return nil, err
	}

	return fact_result, nil
}

func (r factRepositoryDB) Update(fact models.FactTable) (*models.FactTable, error) {

	result := r.db.Model(&fact).Where("id = ?", fact.ID).Updates(fact)
	if result.Error != nil {
		return nil, result.Error
	}

	return &fact, nil
}

func (r factRepositoryDB) Delete(id uint) error {
	fact := models.FactTable{}

	result := r.db.Where("id = ?", id).Delete(&fact)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r factRepositoryDB) GetAll() ([]models.FactTable, error) {
	facts := []models.FactTable{}

	r.db.Find(&facts)

	return facts, nil
}

func (r factRepositoryDB) GetById(id uint) (*models.FactTable, error) {

	fact := models.FactTable{}
	result := r.db.Where("id = ?", id).First(&fact)
	if result.Error != nil {
		return nil, result.Error
	}
	return &fact, nil
}

// Example remove permanently
func (r factRepositoryDB) DeletePermanently(id uint) error {
	fact := models.FactTable{}

	result := r.db.Unscoped().Where("id = ?", id).Delete(&fact)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
