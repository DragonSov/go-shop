package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"lessons/api/presenter"
	"lessons/pkg/common/config"
	"lessons/pkg/entity"
	"lessons/pkg/item"
	"lessons/pkg/user"
)

type ItemHandlers interface {
	CreateItem() fiber.Handler
	GetItems() fiber.Handler
	GetItemByID() fiber.Handler
}

type itemHandlers struct {
	userService user.Service
	itemService item.Service
	cfg         *config.Config
}

func NewItemHandlers(userService user.Service, itemService item.Service, cfg *config.Config) ItemHandlers {
	return &itemHandlers{
		userService: userService,
		itemService: itemService,
		cfg:         cfg,
	}
}

func (h *itemHandlers) CreateItem() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.GetReqHeaders()["Authorization"]
		currentUser, err := getCurrentUser(h.userService, authHeader)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		if currentUser.Role != entity.AdminRole {
			err = NotEnoughPermissions
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		createItem := &entity.ItemCreate{}
		if err = ctx.BodyParser(createItem); err != nil {
			err = BodyParsingError
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		} else if err = createItem.PrepareCreate(); err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		newItem, err := h.itemService.Create(createItem)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		return ctx.JSON(presenter.ItemSuccessResponse(newItem))
	}
}

func (h *itemHandlers) GetItems() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sorting := &entity.Sorting{}
		if err := ctx.QueryParser(sorting); err != nil {
			err = QueryParsingError
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		} else if err = sorting.Validate(); err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		items, err := h.itemService.Get(sorting.Limit, sorting.Offset)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		return ctx.JSON(presenter.ItemsSuccessResponse(items))
	}
}

func (h *itemHandlers) GetItemByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		itemID, err := uuid.Parse(ctx.Params("item_id"))
		if err != nil {
			err = UUIDIncorrect
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		itemData, err := h.itemService.GetByID(itemID)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		return ctx.JSON(presenter.ItemSuccessResponse(itemData))
	}
}
