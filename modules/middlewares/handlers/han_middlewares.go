package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/PeteCroton/go-basic/helpers"
	serv "github.com/PeteCroton/go-basic/modules/middlewares/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middlware-001"
	jwtAuthErr     middlewareHandlersErrCode = "middlware-002"
	paramsCheckErr middlewareHandlersErrCode = "middlware-003"
	authorizeErr   middlewareHandlersErrCode = "middlware-004"
	apiKeyErr      middlewareHandlersErrCode = "middlware-005"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRoleId ...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
	StreamingFile() fiber.Handler
}

type middlewaresHandler struct {
	middlewaresSrv serv.IMiddlewaresService
}

func MiddlewaresHandler(middlewaresSrv serv.IMiddlewaresService) IMiddlewaresHandler {
	return &middlewaresHandler{middlewaresSrv: middlewaresSrv}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return helpers.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"Rotuer not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := helpers.ParseToken(os.Getenv("JWT_SECRET_KEY"))
		if err != nil {
			return helpers.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewaresSrv.FindAccessToken(claims.Id, token) {
			return helpers.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		// Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Locals("userRoleId").(int) == 2 {
			return c.Next()
		}
		if c.Params("user_id") != userId {
			return helpers.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"never gonna give you up",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewaresHandler) Authorize(expectRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return helpers.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErr),
				"user_id is not int type",
			).Res()
		}

		roles, err := h.middlewaresSrv.FindRole()
		if err != nil {
			return helpers.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizeErr),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, v := range expectRoleId {
			sum += v
		}

		expectedValueBinary := helpers.BinaryConverter(sum, len(roles))
		userValueBinary := helpers.BinaryConverter(userRoleId, len(roles))

		// user ->     0 1 0
		// expected -> 1 1 0

		for i := range userValueBinary {
			if userValueBinary[i]&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}
		return helpers.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission to access",
		).Res()
	}
}

func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	// return func(c *fiber.Ctx) error {
	// 	//key := c.Get("X-Api-Key")
	// 	if _, err := helpers.ParseApiKey(os.Getenv("APP_API_KEY")); err != nil {
	// 		return helpers.NewResponse(c).Error(
	// 			fiber.ErrUnauthorized.Code,
	// 			string(apiKeyErr),
	// 			"apikey is invalid or required",
	// 		).Res()
	// 	}
	// 	return c.Next()
	// }

	return nil
}

// Streaming file
func (h *middlewaresHandler) StreamingFile() fiber.Handler {
	return filesystem.New(filesystem.Config{
		Root: http.Dir("./assets/images"),
	})
}
