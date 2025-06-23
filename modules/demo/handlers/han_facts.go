package handlers

import (
	"strconv"

	"github.com/PeteCroton/go-basic/modules/demo/models"
	services "github.com/PeteCroton/go-basic/modules/demo/services"
	"github.com/gofiber/fiber/v2"
)

// Add service methods here
type factHandler struct {
	factSrv services.FactService
}

func NewFactHandler(factSrv services.FactService) factHandler {
	return factHandler{factSrv: factSrv}
}

//=============================
//Applications
//=============================

func (h factHandler) ListFacts(c *fiber.Ctx) error {
	//facts := []models.FactResponse{}
	//var error error
	facts, error := h.factSrv.GetList()

	if error != nil {
		return c.Render("index", fiber.Map{
			"Title":    "Tutorial Question",
			"Subtitle": "Facts for funtimes with friends!",
			"Facts":    nil,
		})
	}
	//configs.DB.Db.Find(&facts)
	return c.Render("index", fiber.Map{
		"Title":    "Tutoaial Question",
		"Subtitle": "Facts for funtimes with friends!",
		"Facts":    facts,
	})
}

func (h factHandler) NewFactView(c *fiber.Ctx) error {
	return c.Render("new", fiber.Map{
		"Title":    "New Fact",
		"Subtitle": "Add a cool fact!",
	})
}

func (h factHandler) CreateFact(c *fiber.Ctx) error {
	fact := new(models.FactTable)
	if err := c.BodyParser(fact); err != nil {
		return h.NewFactView(c)
	}

	// result := configs.DB.Db.Create(&fact)
	// if result.Error != nil {
	// 	return h.NewFactView(c)
	// }

	factRequest := models.FactRequest{
		ID:       fact.ID,
		Question: fact.Question,
		Answer:   fact.Answer,
	}

	factResponse, err := h.factSrv.CrateData(factRequest)
	if err != nil {
		return h.NewFactView(c)
	}
	_ = factResponse
	return h.ListFacts(c)
}

func (h factHandler) ShowFact(c *fiber.Ctx) error {
	// fact := models.Fact{}
	// id := c.Params("id")
	// result := configs.DB.Db.Where("id = ?", id).First(&fact)
	// if result.Error != nil {
	// 	return NotFound(c)
	// }
	// return c.Status(fiber.StatusOK).Render("show", fiber.Map{
	// 	"Title": "Single Fact",
	// 	"Fact":  fact,
	// })

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NotFound(c)
	}
	factResponse, err2 := h.factSrv.GetById(id)
	if err2 != nil {
		return NotFound(c)
	}
	fact := models.FactTable{
		ID:       factResponse.ID,
		Question: factResponse.Question,
		Answer:   factResponse.Answer,
	}

	return c.Status(fiber.StatusOK).Render("show", fiber.Map{
		"Title": "Single Fact",
		"Fact":  fact,
	})

}

func (h factHandler) EditFact(c *fiber.Ctx) error {
	// fact := models.Fact{}
	// id := c.Params("id")

	// result := configs.DB.Db.Where("id = ?", id).First(&fact)
	// if result.Error != nil {
	// 	return NotFound(c)
	// }

	// return c.Render("edit", fiber.Map{
	// 	"Title":    "Edit Fact",
	// 	"Subtitle": "Edit your interesting fact",
	// 	"Fact":     fact,
	// })
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NotFound(c)
	}
	factResponse, err2 := h.factSrv.GetById(id)
	if err2 != nil {
		return NotFound(c)
	}
	fact := models.FactTable{
		ID:       factResponse.ID,
		Question: factResponse.Question,
		Answer:   factResponse.Answer,
	}
	return c.Render("edit", fiber.Map{
		"Title":    "Edit Fact",
		"Subtitle": "Edit your interesting fact",
		"Fact":     fact,
	})
}

func (h factHandler) UpdateFact(c *fiber.Ctx) error {
	// fact := new(models.Fact)
	// id := c.Params("id")

	// if err := c.BodyParser(fact); err != nil {
	// 	return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	// }

	// result := configs.DB.Db.Model(&fact).Where("id = ?", id).Updates(fact)
	// if result.Error != nil {
	// 	return h.EditFact(c)
	// }

	// return h.ShowFact(c)

	fact := new(models.FactTable)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NotFound(c)
	}
	if err := c.BodyParser(fact); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	}

	factRequest := models.FactRequest{
		ID:       uint(id),
		Question: fact.Question,
		Answer:   fact.Answer,
	}

	factResponse, err2 := h.factSrv.UpdateData(factRequest)
	if err2 != nil {
		return NotFound(c)
	}
	_ = factResponse

	return h.ShowFact(c)
}

func (h factHandler) DeleteFact(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NotFound(c)
	}

	err2 := h.factSrv.RemoveData(id)
	if err2 != nil {
		return NotFound(c)
	}

	return h.ListFacts(c)
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendFile("./public/404.html")
}

// =============================
// API Methods
// =============================
func (h factHandler) GetDataByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"success": false,
		})
	}
	fact, err := h.factSrv.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"success": false,
		})
	}
	//return c.Status(fiber.StatusOK).JSON(fact)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    fact,
		"message": "data query successful",
		"success": true,
	})
}
func (h factHandler) GetAll(c *fiber.Ctx) error {

	facts, error := h.factSrv.GetList()
	if error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": error.Error(),
			"success": false,
		})
	}
	//return c.Status(fiber.StatusOK).JSON(facts)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    facts,
		"message": "data query successful",
		"success": true,
	})
}

func (h factHandler) CreateData(c *fiber.Ctx) error {
	fact := new(models.FactRequest)

	if err := c.BodyParser(fact); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"success": false,
		})
	}

	factResponse, err := h.factSrv.CrateData(*fact)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"success": false,
		})
	}

	//return c.Status(fiber.StatusOK).JSON(factResult)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    factResponse,
		"message": "Insert successful",
		"success": true,
	})
}

func (h factHandler) UpdateData(c *fiber.Ctx) error {
	fact := new(models.FactRequest)

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error() + c.Params("id"),
			"success": false,
			"code":    fiber.StatusNotAcceptable,
		})
	}
	fact.ID = uint(id)
	fact.Question = c.FormValue("question")
	fact.Answer = c.FormValue("answer")

	if err := c.BodyParser(fact); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"success": false,
			"code":    fiber.StatusInternalServerError,
		})
	}

	factResponse, err2 := h.factSrv.UpdateData(*fact)
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
		"data":    factResponse,
		"message": "Update successful",
		"success": true,
		"code":    fiber.StatusOK,
	})
}

func (h factHandler) DeleteData(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
			"success": false,
			"code":    fiber.StatusInternalServerError,
		})
	}

	err2 := h.factSrv.RemoveData(id)
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
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
