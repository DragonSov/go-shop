package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"lessons/api/presenter"
	"lessons/pkg/common/config"
	"lessons/pkg/entity"
	"lessons/pkg/item"
	"lessons/pkg/order"
	"lessons/pkg/user"
)

type OrderHandlers interface {
	CreateOrder() fiber.Handler
	GetOrders() fiber.Handler
	GetOrderByID() fiber.Handler
}

type orderHandlers struct {
	userService  user.Service
	itemService  item.Service
	orderService order.Service
	cfg          *config.Config
}

func NewOrderHandlers(userService user.Service, itemService item.Service, orderService order.Service, cfg *config.Config) OrderHandlers {
	return &orderHandlers{
		userService:  userService,
		itemService:  itemService,
		orderService: orderService,
		cfg:          cfg,
	}
}

func (h *orderHandlers) CreateOrder() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.GetReqHeaders()["Authorization"]
		currentUser, err := getCurrentUser(h.userService, authHeader)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.OrderErrorResponse(err))
		}

		if currentUser.Role != entity.AdminRole {
			err = NotEnoughPermissions
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.OrderErrorResponse(err))
		}

		orderData := &entity.OrderCreate{}
		if err = ctx.BodyParser(orderData); err != nil {
			err = BodyParsingError
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.OrderErrorResponse(err))
		}
		createOrder := &entity.Order{
			Items:     orderData.Items,
			OrderedBy: currentUser.ID,
		}

		newOrder, err := h.orderService.Create(createOrder)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.OrderErrorResponse(err))
		}

		return ctx.JSON(presenter.OrderSuccessResponse(newOrder))
	}
}

func (h *orderHandlers) GetOrders() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.GetReqHeaders()["Authorization"]
		currentUser, err := getCurrentUser(h.userService, authHeader)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.OrderErrorResponse(err))
		}

		sorting := &entity.Sorting{}
		if err := ctx.QueryParser(sorting); err != nil {
			err = QueryParsingError
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		} else if err = sorting.Validate(); err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		orders, err := h.orderService.Get(sorting.Limit, sorting.Offset, currentUser.ID)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		return ctx.JSON(presenter.OrdersSuccessResponse(orders))
	}
}

func (h *orderHandlers) GetOrderByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		orderID, err := uuid.Parse(ctx.Params("order_id"))
		if err != nil {
			err = UUIDIncorrect
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.OrderErrorResponse(err))
		}

		orderData, err := h.orderService.GetByID(orderID)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.OrderErrorResponse(err))
		}

		return ctx.JSON(presenter.OrderSuccessResponse(orderData))
	}
}
