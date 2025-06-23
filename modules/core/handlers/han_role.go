package handlers

import (
	"strconv"

	"github.com/PeteCroton/go-basic/helpers"
	"github.com/PeteCroton/go-basic/modules/core/models"
	"github.com/PeteCroton/go-basic/modules/core/services"
	"github.com/gofiber/fiber/v2"
)

// Add service methods here
type roleHandler struct {
	roleSrv services.RoleService
}

func NewRoleHandler(roleSrv services.RoleService) roleHandler {
	return roleHandler{roleSrv: roleSrv}
}

// =============================
// APIs
// =============================
const (
	dataInsertError    string = "role-001"
	dataUpdateError    string = "role-002"
	dateDeleteError    string = "role-003"
	dataQueryError     string = "role-004"
	dataBodyParseError string = "role-005"
)

func (h roleHandler) GetAll(c *fiber.Ctx) error {

	roles, error := h.roleSrv.GetAll()
	if error != nil {
		return helpers.NewResponse(c).Error(fiber.ErrBadRequest.Code, dataQueryError, error.Error()).Res()
	}
	//helpers return
	return helpers.NewResponse(c).Success(fiber.StatusOK, roles).Res()

}

func (h roleHandler) GetById(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.NewResponse(c).Error(fiber.ErrBadRequest.Code, dataQueryError, err.Error()).Res()
	}
	role, err := h.roleSrv.GetById(uint(id))
	if err != nil {
		return helpers.NewResponse(c).Error(fiber.ErrBadRequest.Code, dataQueryError, err.Error()).Res()
	}
	return helpers.NewResponse(c).Success(fiber.StatusOK, role).Res()
}

func (h roleHandler) CrateData(c *fiber.Ctx) error {

	role := new(models.RoleTable)

	if err := c.BodyParser(role); err != nil {
		return helpers.NewResponse(c).Error(fiber.ErrBadRequest.Code, dataBodyParseError, err.Error()).Res()
	}
	dataResponse, err := h.roleSrv.CrateData(*role)
	if err != nil {
		return helpers.NewResponse(c).Error(fiber.ErrBadRequest.Code, dataInsertError, err.Error()).Res()
	}
	return helpers.NewResponse(c).Success(fiber.StatusOK, dataResponse).Res()

}

func (h roleHandler) UpdateData(c *fiber.Ctx) error {
	data := new(models.RoleTable)

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error() + c.FormValue("id"),
			"success": false,
			"code":    fiber.StatusNotAcceptable,
		})
	}
	data.ID = uint(id)
	data.Title = c.FormValue("title")

	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"success": false,
			"code":    fiber.StatusInternalServerError,
		})
	}

	dataResponse, err2 := h.roleSrv.UpdateData(*data)
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err2.Error(),
			"success": false,
			"code":    fiber.StatusInternalServerError,
		})
	}

	//return c.Status(fiber.StatusOK).JSON(factResponse)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    dataResponse,
		"message": "Update successful",
		"success": true,
		"code":    fiber.StatusOK,
	})
}

func (h roleHandler) DeleteData(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"success": false,
			"code":    fiber.StatusInternalServerError,
		})
	}

	err2 := h.roleSrv.DeleteData(uint(id))
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err2.Error(),
			"success": false,
			"code":    fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    nil,
		"message": "Data Deleted",
		"success": true,
		"code":    fiber.StatusOK,
	})
}
