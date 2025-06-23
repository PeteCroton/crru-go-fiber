package handlers

import (
	"strconv"

	"github.com/PeteCroton/go-basic/modules/core/models"
	"github.com/PeteCroton/go-basic/modules/core/services"
	"github.com/gofiber/fiber/v2"
)

// Add service methods here
// Define แล้วใน han_user.go
// type userHandler struct {
// 	userSrv services.UserService
// }

func NewUserLoginHandler(userSrv services.UserService) userHandler {
	return userHandler{userSrv: userSrv}
}

//=============================
// APIs
//=============================

func (h userHandler) Login(c *fiber.Ctx) error {
	data := new(models.LoginRequest)

	data.Username = c.FormValue("login_username")
	data.Password = c.FormValue("login_password")

	token, err := h.userSrv.FindByCredentials(*data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"code":    fiber.StatusInternalServerError,
			"success": false,
		})
	}

	user_id, errx := strconv.Atoi(token.Id)
	if errx != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data":          nil,
			"access_token":  nil,
			"refresh_token": nil,
			"message":       "Not found user ID",
			"code":          fiber.StatusNotFound,
			"success":       false,
		})
	}

	access_token := token.AccessToken
	refresh_token := token.RefreshToken

	rs_data, err := h.userSrv.GetById(uint(user_id))
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data":          nil,
			"access_token":  nil,
			"refresh_token": nil,
			"message":       "Login Fail",
			"code":          fiber.StatusNotFound,
			"success":       false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_data": fiber.Map{
			"UserID":     rs_data.ID,
			"UserRoleID": rs_data.Role_ID,
			"Email":      rs_data.Email,
		},
		"access_token":  access_token,
		"refresh_token": refresh_token,
		"message":       "Login Success",
		"code":          fiber.StatusOK,
		"success":       true,
	})
}
