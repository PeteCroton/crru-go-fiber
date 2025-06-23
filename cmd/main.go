package main

import (
	"github.com/PeteCroton/go-basic/configs"
	module_core "github.com/PeteCroton/go-basic/modules/core/routes"
	"github.com/PeteCroton/go-basic/modules/demo/handlers"
	module_demo "github.com/PeteCroton/go-basic/modules/demo/routes"

	hand_mid "github.com/PeteCroton/go-basic/modules/middlewares/handlers"
	repo_mid "github.com/PeteCroton/go-basic/modules/middlewares/repo"
	serv_mid "github.com/PeteCroton/go-basic/modules/middlewares/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {

	configs.ConnectDb()

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		ServerHeader: "Go Basic by CrotonDev",
		Views:        engine,
		ViewsLayout:  "layouts/main",
		Prefork:      false,
	})

	basic_user := map[string]string{
		"qkgb4dTRrBSpaC7M": "LeDSk6tx5f58rB5m2mLYcWaXLuADU7CY",
		"admin":            "123456",
	}

	//Middlewares
	db := configs.DbConnect()
	defer db.Close()

	middleRepositoryDB := repo_mid.MiddlewaresRepository(db)
	middleService := serv_mid.MiddlewaresService(middleRepositoryDB)
	middleHandler := hand_mid.MiddlewaresHandler(middleService)

	app.Use(middleHandler.Cors())
	app.Use(middleHandler.Logger())
	app.Use(middleHandler.StreamingFile())

	// module routes of api
	module_demo.SetupRoutesDemo(app)
	module_core.SetupRoutesCore(app, basic_user)

	app.Static("/", "./public")
	//app.Use(middleHandler.RouterCheck())
	app.Use(handlers.NotFound)

	app.Listen(":3000")
}
