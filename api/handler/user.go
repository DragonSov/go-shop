package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"lessons/api/presenter"
	"lessons/pkg/common/config"
	"lessons/pkg/common/utils"
	"lessons/pkg/entity"
	"lessons/pkg/user"
)

type UserHandlers interface {
	CreateUser() fiber.Handler
	LoginUser() fiber.Handler
	GetCurrentUser() fiber.Handler
	GetUserByID() fiber.Handler
	UpdateCurrentUser() fiber.Handler
}

type userHandlers struct {
	userService user.Service
	cfg         *config.Config
}

func NewUserHandlers(userService user.Service, cfg *config.Config) UserHandlers {
	return &userHandlers{
		userService: userService,
		cfg:         cfg,
	}
}

func (h *userHandlers) CreateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		data := entity.UserCredentials{}
		if err := ctx.BodyParser(&data); err != nil {
			err = BodyParsingError
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		createUser := &entity.User{
			Login:    data.Login,
			Password: data.Password,
		}
		if err := createUser.PrepareCreate(); err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		newUser, err := h.userService.Create(createUser)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		return ctx.JSON(presenter.UserSuccessResponse(newUser))
	}
}

func (h *userHandlers) LoginUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userCreds := entity.UserCredentials{}
		if err := ctx.BodyParser(&userCreds); err != nil {
			err = BodyParsingError
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		} else if err = userCreds.ToLowerCase(); err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		selectedUser, err := h.userService.GetByLogin(userCreds.Login)
		if selectedUser == nil || selectedUser.ComparePassword(userCreds.Password) != nil {
			err = LoginOrPasswordIncorrect
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		} else if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		userToken, err := utils.GenerateJWTToken(selectedUser, h.cfg)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		return ctx.JSON(presenter.TokenSuccessResponse(selectedUser.ID, userToken))
	}
}

func (h *userHandlers) GetCurrentUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.GetReqHeaders()["Authorization"]
		currentUser, err := getCurrentUser(h.userService, authHeader)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		return ctx.JSON(presenter.UserSuccessResponse(currentUser))
	}
}

func (h *userHandlers) GetUserByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID, err := uuid.Parse(ctx.Params("user_id"))
		if err != nil {
			err = UUIDIncorrect
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		selectedUser, err := h.userService.GetByID(userID)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		return ctx.JSON(presenter.UserSuccessResponse(selectedUser))
	}
}

func (h *userHandlers) UpdateCurrentUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.GetReqHeaders()["Authorization"]
		currentUser, err := getCurrentUser(h.userService, authHeader)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		var data struct {
			Password string `json:"password"`
		}
		if err = ctx.BodyParser(&data); err != nil {
			err = BodyParsingError
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		updatedUser, err := h.userService.UpdatePasswordByID(currentUser.ID, data.Password)
		if err != nil {
			err, statusCode := exceptionHandler(err)
			return ctx.Status(statusCode).JSON(presenter.ItemErrorResponse(err))
		}

		return ctx.JSON(presenter.UserSuccessResponse(updatedUser))
	}
}
