package routes

import (
	"github.com/PeteCroton/go-basic/configs"
	"github.com/PeteCroton/go-basic/modules/demo/handlers"
	repo "github.com/PeteCroton/go-basic/modules/demo/repo/database"
	"github.com/PeteCroton/go-basic/modules/demo/services"
	"github.com/gofiber/fiber/v2"
)


func SetupRoutesDemo(app *fiber.App) {

	factRepositoryDB := repo.NewFactRepositoryDB(configs.ConnectDb())
	factService := services.NewFactService(factRepositoryDB)
	factHandler := handlers.NewFactHandler(factService)

	app.Get("/fact/home", factHandler.ListFacts)
	app.Get("/fact", factHandler.NewFactView)
	app.Post("/fact", factHandler.CreateFact)
	app.Get("/fact/:id", factHandler.ShowFact)
	app.Get("/fact/:id/edit", factHandler.EditFact)
	app.Patch("/fact/:id", factHandler.UpdateFact)
	app.Delete("/fact/:id", factHandler.DeleteFact)

	v1 := app.Group("/api/v1")
	v1.Get("/fact/get_list_data", factHandler.GetAll)
	v1.Get("/fact/get_data/:id", factHandler.GetDataByID) //curl http://127.0.0.1:3000/api/v1/fact/get_data/1
	v1.Post("/fact/create_data", factHandler.CreateData)  //curl -X POST -H "Content-Type: application/json" -d '{"Question":"Hello Wortld", "Answer":"xxxx"}' http://localhost:3000/api/v1/fact/create_data
	v1.Patch("/fact/update_data", factHandler.UpdateData)
	v1.Delete("/fact/delete_data", factHandler.DeleteData)
}
