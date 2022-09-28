package presenter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"lessons/pkg/entity"
	"time"
	"unicode"
)

// User это данные, которые можно отдавать пользователю с помощью Response
type User struct {
	ID        uuid.UUID       `json:"id"`
	Login     string          `json:"login"`
	Role      entity.UserRole `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
}

type Token struct {
	UserID uuid.UUID `json:"user_id"`
	Token  string    `json:"token"`
}

// UserSuccessResponse это SuccessResponse для хендлера
func UserSuccessResponse(data *entity.User) *fiber.Map {
	user := User{
		ID:        data.ID,
		Login:     data.Login,
		Role:      data.Role,
		CreatedAt: data.CreatedAt,
	}

	return &fiber.Map{
		"status": true,
		"data":   user,
		"error":  nil,
	}
}

func TokenSuccessResponse(userID uuid.UUID, token string) *fiber.Map {
	data := Token{
		UserID: userID,
		Token:  token,
	}

	return &fiber.Map{
		"status": true,
		"data":   data,
		"error":  nil,
	}
}

// UserErrorResponse это ErrorResponse для хендлера
func UserErrorResponse(err error) *fiber.Map {
	errText := []rune(err.Error())
	errText[0] = unicode.ToUpper(errText[0])

	return &fiber.Map{
		"status": false,
		"data":   nil,
		"error":  string(errText),
	}
}
