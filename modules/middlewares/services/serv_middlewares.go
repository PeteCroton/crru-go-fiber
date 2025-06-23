package services

import (
	"github.com/PeteCroton/go-basic/modules/middlewares/models"
	repo "github.com/PeteCroton/go-basic/modules/middlewares/repo"
)

type IMiddlewaresService interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*models.Role, error)
}

type middlewaresService struct {
	middlewaresRepository repo.IMiddlewaresRepository
}

func MiddlewaresService(middlewaresRepository repo.IMiddlewaresRepository) IMiddlewaresService {
	return &middlewaresService{
		middlewaresRepository: middlewaresRepository,
	}
}

func (u *middlewaresService) FindAccessToken(userId, accessToken string) bool {
	return u.middlewaresRepository.FindAccessToken(userId, accessToken)
}

func (u *middlewaresService) FindRole() ([]*models.Role, error) {
	roles, err := u.middlewaresRepository.FindRole()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
