package services

import (
	"strconv"

	helpers "github.com/PeteCroton/go-basic/helpers"
	models "github.com/PeteCroton/go-basic/modules/core/models"
	repo "github.com/PeteCroton/go-basic/modules/core/repo/database"
	modelMiddel "github.com/PeteCroton/go-basic/modules/middlewares/models"
	"github.com/go-playground/validator/v10"
	//jtoken "github.com/golang-jwt/jwt/v4"
)

type UserService interface {
	CrateData(models.UserTable) (*models.UserTable, error)
	UpdateData(models.UserTable) (*models.UserTable, error)
	DeleteData(uint) error
	GetAll() ([]models.UserTable, error)
	GetById(uint) (*models.UserTable, error)
	FindByCredentials(models.LoginRequest) (modelMiddel.UserToken, error)
}

type userService struct {
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) UserService {
	return userService{userRepo: userRepo}
}

func (s userService) CrateData(data models.UserTable) (*models.UserTable, error) {
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

	dataResponse, err := s.userRepo.Create(data)
	if err != nil {
		return nil, err
	}
	return dataResponse, nil
}

func (s userService) UpdateData(data models.UserTable) (*models.UserTable, error) {
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
	dataResponse, err := s.userRepo.Update(data)
	if err != nil {
		return nil, err
	}

	return dataResponse, nil
}

func (s userService) DeleteData(id uint) error {

	if id == 0 {
		return helpers.NewValidationError("ID is not valid")
	}

	err := s.userRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s userService) GetAll() ([]models.UserTable, error) {

	listData, err := s.userRepo.GetAll()

	if err != nil {
		return nil, err
	}

	return listData, nil
}

func (s userService) GetById(id uint) (*models.UserTable, error) {

	//Validate
	if id == 0 {
		return nil, helpers.NewValidationError("ID is not valid")
	}

	dataDB, err := s.userRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return dataDB, nil
}

func (s userService) FindByCredentials(data models.LoginRequest) (modelMiddel.UserToken, error) {

	//validation
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {

		error_message := ""
		for _, e := range err.(validator.ValidationErrors) {
			error_message += e.Error() + "\n"
		}
		return modelMiddel.UserToken{}, helpers.NewValidationError(error_message)
	}

	dataDB, err := s.userRepo.FindByCredentials(data)
	if err != nil {
		return modelMiddel.UserToken{}, helpers.NewValidationError(err.Error())
	}
	// day := time.Hour * 24
	clam_id := strconv.FormatUint(uint64(dataDB.ID), 10)
	clam_role_id := int(dataDB.Role_ID)

	userClaims := &modelMiddel.UserClaims{
		Id:     clam_id,
		RoleId: clam_role_id,
	}
	// สร้าง access token
	access, err := helpers.NewGobasicAuth(helpers.Access, userClaims)
	if err != nil {
		return modelMiddel.UserToken{}, helpers.NewValidationError(err.Error())
	}

	// สร้าง refresh token
	refresh, err := helpers.NewGobasicAuth(helpers.Refresh, userClaims)
	if err != nil {
		return modelMiddel.UserToken{}, helpers.NewValidationError(err.Error())
	}

	// สร้าง token struct
	userToken := modelMiddel.UserToken{
		Id:           strconv.FormatUint(uint64(dataDB.ID), 10),
		AccessToken:  access.SignToken(),
		RefreshToken: refresh.SignToken(),
	}

	return userToken, nil

	// day := time.Hour * 24
	// // Create the JWT claims, which includes the user ID and expiry time
	// claims := jtoken.MapClaims{
	// 	"ID":      dataDB.ID,
	// 	"email":   dataDB.Email,
	// 	"role_id": dataDB.Role_ID,
	// 	"message": "This is a custom claim",
	// 	"exp":     time.Now().Add(day * 1).Unix(),
	// }
	// Create token
	//token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	// t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	// if err != nil {

	// 	return "", helpers.NewValidationError(err.Error())
	// }
	// return t, nil

}
