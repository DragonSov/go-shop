package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
	"lessons/api/route"
	"lessons/pkg/common/config"
	"lessons/pkg/common/database"
	"lessons/pkg/item"
	"lessons/pkg/order"
	"lessons/pkg/user"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.ConnectDatabase(cfg.DSN)
	if err != nil {
		panic(err)
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, cfg)

	itemRepo := item.NewRepository(db)
	itemService := item.NewService(itemRepo)

	orderRepo := order.NewRepository(db)
	orderService := order.NewService(itemRepo, orderRepo)

	app := fiber.New()

	app.Use(logger.New())

	route.UserRouter(app, userService, cfg)
	route.ItemRouter(app, userService, itemService, cfg)
	route.OrderRouter(app, userService, itemService, orderService, cfg)

	panic(app.Listen(cfg.Addr))
}
