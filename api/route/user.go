package route

import (
	"github.com/gofiber/fiber/v2"
	"lessons/api/handler"
	"lessons/pkg/common/config"
	"lessons/pkg/user"
)

func UserRouter(app *fiber.App, userService user.Service, cfg *config.Config) {
	group := app.Group("/users")
	handlers := handler.NewUserHandlers(userService, cfg)

	group.Post("/", handlers.CreateUser())
	group.Post("/login", handlers.LoginUser())
	group.Get("/me", handlers.GetCurrentUser())
	group.Get("/:user_id", handlers.GetUserByID())
	group.Patch("/", handlers.UpdateCurrentUser())
}
