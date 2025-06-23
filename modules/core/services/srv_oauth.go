package services

import (
	"github.com/PeteCroton/go-basic/helpers"

	models "github.com/PeteCroton/go-basic/modules/core/models"
	repo "github.com/PeteCroton/go-basic/modules/core/repo/database"
)

type OauthService interface {
	CrateData(models.OauthTable) (*models.OauthResponse, error)
	UpdateData(models.OauthTable) (*models.OauthResponse, error)
	GetList() ([]models.OauthResponse, error)
	GetById(int) (*models.OauthResponse, error)
	RemoveData(int) error
}

type oauthService struct {
	oauthRepo repo.OauthRepository
}

func NewOuthService(oauthRepo repo.OauthRepository) OauthService {
	return oauthService{oauthRepo: oauthRepo}
}

func (s oauthService) CrateData(data models.OauthTable) (*models.OauthResponse, error) {
	//Validate
	if data.AccessToken == "" {
		return nil, helpers.NewValidationError("Data is not valid")
	}

	//Request Model to Gorm Model
	//dataTable := requestModelToDbModel(data)
	dataResponse, err := s.oauthRepo.Create(data)
	if err != nil {
		return nil, err
	}

	//Gorm Model to Response Model
	result := dbmodelToResponse(*dataResponse)

	return &result, nil
}

func (s oauthService) UpdateData(data models.OauthTable) (*models.OauthResponse, error) {
	// //Validate
	// if fact.Question == "" || fact.Answer == "" {
	// 	return nil, helpers.NewValidationError("Data is not valid")
	// }
	//Request Model to Gorm Model
	//dataDB := requestModelToDbModel(data)

	factResponse, err := s.oauthRepo.Update(data)
	if err != nil {
		return nil, err
	}

	//Gorm Model to Response Model
	result := dbmodelToResponse(*factResponse)

	return &result, nil
}

func (s oauthService) GetList() ([]models.OauthResponse, error) {

	listData, err := s.oauthRepo.GetAll()

	if err != nil {
		return nil, err
	}

	listResponse := []models.OauthResponse{}

	for _, row := range listData {
		//Gorm Model to Response Model
		result := dbmodelToResponse(row)
		listResponse = append(listResponse, result)
	}

	return listResponse, nil
}

func (s oauthService) GetById(id int) (*models.OauthResponse, error) {
	//Validate
	if id == 0 {
		return nil, helpers.NewValidationError("ID is not valid")
	}

	dataDB, err := s.oauthRepo.GetById(uint(id))

	if err != nil {
		return nil, err
	}

	//Gorm Model to Response Model
	result := dbmodelToResponse(*dataDB)

	return &result, nil
}

func (s oauthService) RemoveData(id int) error {
	//Validate
	if id == 0 {
		return helpers.NewValidationError("ID is not valid")
	}

	err := s.oauthRepo.Delete(uint(id))
	if err != nil {
		return err
	}

	return nil
}

// =====================================================================
func dbmodelToResponse(data models.OauthTable) models.OauthResponse {
	return models.OauthResponse{
		ID:           data.ID,
		UserID:       data.UserID,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}
}

func requestModelToDbModel(data models.OauthResponse) models.OauthTable {
	return models.OauthTable{
		ID:           data.ID,
		UserID:       data.UserID,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}
}
