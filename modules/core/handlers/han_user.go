package handlers

import (
	"strconv"
	"time"

	"github.com/PeteCroton/go-basic/modules/core/models"
	"github.com/PeteCroton/go-basic/modules/core/services"
	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
)

// Add service methods here
type userHandler struct {
	userSrv services.UserService
}

func NewUserHandler(userSrv services.UserService) userHandler {
	return userHandler{userSrv: userSrv}
}

//=============================
// APIs
//=============================

func (h userHandler) GetAll(c *fiber.Ctx) error {

	listData, error := h.userSrv.GetAll()
	if error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": error.Error(),
			"code":    fiber.StatusInternalServerError,
			"success": false,
		})
	}
	//return c.Status(fiber.StatusOK).JSON(roles)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    listData,
		"message": "data query successful",
		"code":    fiber.StatusOK,
		"success": true,
	})
}

func (h userHandler) GetById(c *fiber.Ctx) error {
	//id := c.Params("id")
	id, err := strconv.Atoi(c.Params("id")) //Convert string to int
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"code":    fiber.StatusUnprocessableEntity,
			"success": false,
		})
	}
	role, err := h.userSrv.GetById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"code":    fiber.StatusUnprocessableEntity,
			"success": false,
		})
	}
	//return c.Status(fiber.StatusOK).JSON(role)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    role,
		"message": "data query successful",
		"code":    fiber.StatusOK,
		"success": true,
	})

}

func (h userHandler) CrateData(c *fiber.Ctx) error {

	data := new(models.UserTable)

	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"code":    fiber.StatusInternalServerError,
			"success": false,
		})
	}
	dataResponse, err := h.userSrv.CrateData(*data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"code":    fiber.StatusInternalServerError,
			"success": false,
		})
	}

	//return c.Status(fiber.StatusOK).JSON(factResult)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    dataResponse,
		"message": "Insert successful",
		"code":    fiber.StatusOK,
		"success": true,
	})

}

func (h userHandler) UpdateData(c *fiber.Ctx) error {
	data := new(models.UserTable)

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error() + c.FormValue("id"),
			"code":    fiber.StatusNotAcceptable,
			"success": false,
		})
	}
	data.ID = uint(id)
	// data.ID = c.Params("id")
	user_id, err := strconv.Atoi(c.FormValue("role_id"))
	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error() + c.FormValue("role_id"),
			"code":    fiber.StatusNotAcceptable,
			"success": false,
		})
	}
	data.Role_ID = uint(user_id)
	data.Username = c.FormValue("username")
	data.Password = c.FormValue("password")
	data.Email = c.FormValue("email")

	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"code":    fiber.StatusInternalServerError,
			"success": false,
		})
	}

	dataResponse, err2 := h.userSrv.UpdateData(*data)
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err2.Error(),
			"code":    fiber.StatusInternalServerError,
			"success": false,
		})
	}

	//return c.Status(fiber.StatusOK).JSON(factResponse)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    dataResponse,
		"message": "Update successful",
		"code":    fiber.StatusOK,
		"success": true,
	})
}

func (h userHandler) DeleteData(c *fiber.Ctx) error {
	//id := c.Params("id")

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"code":    fiber.StatusInternalServerError,
			"success": false,
		})
	}

	err2 := h.userSrv.DeleteData(uint(id))
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err2.Error(),
			"code":    fiber.StatusInternalServerError,
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    nil,
		"message": "Data Deleted",
		"code":    fiber.StatusOK,
		"success": true,
	})
}

// Protected route
func (h userHandler) Protected(c *fiber.Ctx) error {
	// Get the user from the context and return it
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	email := claims["email"].(string)
	message := claims["message"].(string)

	expiration := int64(claims["exp"].(float64))
	if time.Now().Unix() > expiration {
		return c.SendString("Token has expired")
	}

	return c.SendString("Welcome 👋 :" + email + " msg:" + message)
}
