package presenter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"lessons/pkg/entity"
	"time"
	"unicode"
)

type Order struct {
	ID        uuid.UUID         `json:"id"`
	Items     map[uuid.UUID]int `json:"items"`
	Price     float32           `json:"price"`
	OrderedBy uuid.UUID         `json:"ordered_by"`
	CreatedAt time.Time         `json:"created_at"`
}

func OrderSuccessResponse(data *entity.Order) *fiber.Map {
	order := Order{
		ID:        data.ID,
		Items:     data.Items,
		Price:     data.Price,
		OrderedBy: data.OrderedBy,
		CreatedAt: data.CreatedAt,
	}

	return &fiber.Map{
		"status": true,
		"data":   order,
		"error":  nil,
	}
}

func OrdersSuccessResponse(data *[]entity.Order) *fiber.Map {
	var orders []Order
	for _, orderData := range *data {
		order := Order{
			ID:        orderData.ID,
			Items:     orderData.Items,
			Price:     orderData.Price,
			OrderedBy: orderData.OrderedBy,
			CreatedAt: orderData.CreatedAt,
		}

		orders = append(orders, order)
	}

	return &fiber.Map{
		"status": true,
		"data":   orders,
		"error":  nil,
	}
}

func OrderErrorResponse(err error) *fiber.Map {
	errText := []rune(err.Error())
	errText[0] = unicode.ToUpper(errText[0])

	return &fiber.Map{
		"status": false,
		"data":   nil,
		"error":  string(errText),
	}
}
