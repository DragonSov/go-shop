package route

import (
	"github.com/gofiber/fiber/v2"
	"lessons/api/handler"
	"lessons/pkg/common/config"
	"lessons/pkg/item"
	"lessons/pkg/user"
)

func ItemRouter(app *fiber.App, userService user.Service, itemService item.Service, cfg *config.Config) {
	group := app.Group("/items")
	itemHandlers := handler.NewItemHandlers(userService, itemService, cfg)

	group.Post("/", itemHandlers.CreateItem())
	group.Get("/", itemHandlers.GetItems())
	group.Get("/:item_id", itemHandlers.GetItemByID())
}
