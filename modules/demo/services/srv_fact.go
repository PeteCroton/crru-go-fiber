package services

import (
	errs "github.com/PeteCroton/go-basic/helpers"
	models "github.com/PeteCroton/go-basic/modules/demo/models"
	repo "github.com/PeteCroton/go-basic/modules/demo/repo/database"
)

type FactService interface {
	CrateData(models.FactRequest) (*models.FactResponse, error)
	UpdateData(models.FactRequest) (*models.FactResponse, error)
	GetList() ([]models.FactResponse, error)
	GetById(int) (*models.FactResponse, error)
	RemoveData(int) error
}


type factService struct {
	factRepo repo.FactRepository
}

func NewFactService(factRepo repo.FactRepository) FactService {
	return factService{factRepo: factRepo}
}

func (s factService) CrateData(fact models.FactRequest) (*models.FactResponse, error) {
	//Validate
	if fact.Question == "" || fact.Answer == "" {
		return nil, errs.NewValidationError("Data is not valid")
	}

	//Request Model to Gorm Model
	factDB := requestModelToDbModel(fact)
	factResponse, err := s.factRepo.Create(factDB)
	if err != nil {
		return nil, err
	}

	//Gorm Model to Response Model
	result := dbmodelToResponse(*factResponse)

	return &result, nil
}

func (s factService) UpdateData(fact models.FactRequest) (*models.FactResponse, error) {
	//Validate
	if fact.Question == "" || fact.Answer == "" {
		return nil, errs.NewValidationError("Data is not valid")
	}
	//Request Model to Gorm Model
	factDB := requestModelToDbModel(fact)

	factResponse, err := s.factRepo.Update(factDB)
	if err != nil {
		return nil, err
	}

	//Gorm Model to Response Model
	result := dbmodelToResponse(*factResponse)

	return &result, nil
}

func (s factService) GetList() ([]models.FactResponse, error) {

	factDBs, err := s.factRepo.GetAll()

	if err != nil {
		return nil, err
	}

	factResponse := []models.FactResponse{}

	for _, factDB := range factDBs {
		//Gorm Model to Response Model
		result := dbmodelToResponse(factDB)
		factResponse = append(factResponse, result)
	}

	return factResponse, nil
}

func (s factService) GetById(id int) (*models.FactResponse, error) {
	//Validate
	if id == 0 {
		return nil, errs.NewValidationError("ID is not valid")
	}

	factDB, err := s.factRepo.GetById(uint(id))

	if err != nil {
		return nil, err
	}

	//Gorm Model to Response Model
	result := dbmodelToResponse(*factDB)

	return &result, nil
}

func (s factService) RemoveData(id int) error {
	//Validate
	if id == 0 {
		return errs.NewValidationError("ID is not valid")
	}

	err := s.factRepo.Delete(uint(id))
	if err != nil {
		return err
	}

	return nil
}

// =====================================================================
func dbmodelToResponse(fact models.FactTable) models.FactResponse {
	return models.FactResponse{
		ID:       fact.ID,
		Question: fact.Question,
		Answer:   fact.Answer,
	}
}

func requestModelToDbModel(fact models.FactRequest) models.FactTable {
	return models.FactTable{
		ID:       fact.ID,
		Question: fact.Question,
		Answer:   fact.Answer,
	}
}
