package presenter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"lessons/pkg/entity"
	"time"
	"unicode"
)

type Item struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

func ItemSuccessResponse(data *entity.Item) *fiber.Map {
	item := Item{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		CreatedAt:   data.CreatedAt,
	}

	return &fiber.Map{
		"status": true,
		"data":   item,
		"error":  nil,
	}
}

func ItemsSuccessResponse(data *[]entity.Item) *fiber.Map {
	var items []Item
	for _, itemData := range *data {
		item := Item{
			ID:          itemData.ID,
			Name:        itemData.Name,
			Description: itemData.Description,
			Price:       itemData.Price,
			CreatedAt:   itemData.CreatedAt,
		}

		items = append(items, item)
	}

	return &fiber.Map{
		"status": true,
		"data":   items,
		"error":  nil,
	}
}

func ItemErrorResponse(err error) *fiber.Map {
	errText := []rune(err.Error())
	errText[0] = unicode.ToUpper(errText[0])

	return &fiber.Map{
		"status": false,
		"data":   nil,
		"error":  string(errText),
	}
}
