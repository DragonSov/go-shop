package route

import (
	"github.com/gofiber/fiber/v2"
	"lessons/api/handler"
	"lessons/pkg/common/config"
	"lessons/pkg/item"
	"lessons/pkg/order"
	"lessons/pkg/user"
)

func OrderRouter(app *fiber.App, userService user.Service, itemService item.Service, orderService order.Service, cfg *config.Config) {
	group := app.Group("/orders")
	orderHandlers := handler.NewOrderHandlers(userService, itemService, orderService, cfg)

	group.Post("/", orderHandlers.CreateOrder())
	group.Get("/", orderHandlers.GetOrders())
	group.Get("/:order_id", orderHandlers.GetOrderByID())
}
