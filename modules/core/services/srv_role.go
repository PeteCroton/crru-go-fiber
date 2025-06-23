package services

import (
	helpers "github.com/PeteCroton/go-basic/helpers"
	models "github.com/PeteCroton/go-basic/modules/core/models"
	repo "github.com/PeteCroton/go-basic/modules/core/repo/database"
	"github.com/go-playground/validator/v10"
)

type RoleService interface {
	CrateData(models.RoleTable) (*models.RoleTable, error)
	UpdateData(models.RoleTable) (*models.RoleTable, error)
	DeleteData(uint) error
	GetAll() ([]models.RoleTable, error)
	GetById(uint) (*models.RoleTable, error)
}

type roleService struct {
	roleRepo repo.RoleRepository
}

func NewRoleService(roleRepo repo.RoleRepository) RoleService {
	return roleService{roleRepo: roleRepo}
}

func (s roleService) CrateData(data models.RoleTable) (*models.RoleTable, error) {
	//validation
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {

		error_message := ""
		for _, e := range err.(validator.ValidationErrors) {
			error_message += e.Error() + "\n"
		}
		return nil, helpers.NewValidationError(error_message)
	}

	dataResponse, err := s.roleRepo.Create(data)
	if err != nil {
		return nil, err
	}
	return dataResponse, nil
}

func (s roleService) UpdateData(data models.RoleTable) (*models.RoleTable, error) {

	if data.Title == "" {
		return nil, helpers.NewValidationError("Data is not valid")
	}
	dataResponse, err := s.roleRepo.Update(data)
	if err != nil {
		return nil, err
	}

	return dataResponse, nil
}

func (s roleService) DeleteData(id uint) error {

	if id == 0 {
		return helpers.NewValidationError("ID is not valid")
	}

	err := s.roleRepo.Delete(uint(id))
	if err != nil {
		return err
	}

	return nil
}

func (s roleService) GetAll() ([]models.RoleTable, error) {
	listData, err := s.roleRepo.GetAll()

	if err != nil {
		return nil, err
	}

	return listData, nil
}

func (s roleService) GetById(id uint) (*models.RoleTable, error) {
	//Validate
	if id == 0 {
		return nil, helpers.NewValidationError("ID is not valid")
	}

	dataDB, err := s.roleRepo.GetById(uint(id))

	if err != nil {
		return nil, err
	}

	return dataDB, nil

}
