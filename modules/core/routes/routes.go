package routes

import (
	"os"

	"github.com/PeteCroton/go-basic/configs"
	"github.com/PeteCroton/go-basic/helpers"
	hand_core "github.com/PeteCroton/go-basic/modules/core/handlers"
	repo_core "github.com/PeteCroton/go-basic/modules/core/repo/database"
	serv_core "github.com/PeteCroton/go-basic/modules/core/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func SetupRoutesCore(app *fiber.App, user_basic map[string]string) {

	roleRepositoryDB := repo_core.NewRoleRepositoryDB(configs.ConnectDb())
	roleService := serv_core.NewRoleService(roleRepositoryDB)
	roleHandler := hand_core.NewRoleHandler(roleService)

	userRepositoryDB := repo_core.NewUserRepositoryDB(configs.ConnectDb())
	userService := serv_core.NewUserService(userRepositoryDB)
	userHandler := hand_core.NewUserHandler(userService)
	userLoginHandler := hand_core.NewUserLoginHandler(userService)

	config := basicauth.Config{
		Users: user_basic,
	}
	basic_auth1 := basicauth.New(config)

	jwt := helpers.NewAuthMiddleware(os.Getenv("JWT_SECRET_KEY"))

	v1 := app.Group("/api/v1/core")
	//role - jwt authen all prefix with /role
	v1.Use("/role", jwt)
	v1.Get("/role/get_list_data", roleHandler.GetAll)
	v1.Get("/role/get_data/:id", roleHandler.GetById)
	v1.Post("/role/create_data", roleHandler.CrateData)
	v1.Patch("/role/update_data", roleHandler.UpdateData)
	v1.Delete("/role/remove_data", roleHandler.DeleteData)

	v1.Use("Test", basic_auth1)
	//user - basic authen
	// v1.Post("/user/create_data", basic_auth1, userHandler.CrateData)
	// v1.Patch("/user/update_data", basic_auth1, userHandler.UpdateData)
	v1.Post("/userlogin/login", userLoginHandler.Login)

	//User
	v1.Get("/user/get_data/:id", jwt, userHandler.GetById)
	v1.Get("/user/get_list_data", jwt, userHandler.GetAll)
	v1.Post("/user/create_data", jwt, userHandler.CrateData)
	v1.Patch("/user/update_data", jwt, userHandler.UpdateData)
	v1.Get("/user/test_jwt", jwt, userHandler.Protected)
	v1.Delete("/user/remove_data", jwt, userHandler.DeleteData)
}
