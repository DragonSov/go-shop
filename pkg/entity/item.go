package entity

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

var (
	NameLengthIncorrect        = errors.New("имя предмета должно быть длиной от 3 до 128 символов")
	DescriptionLengthIncorrect = errors.New("описание предмета должно быть длиной до 1024 символов")
)

type ItemCreate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type Item struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Price       float32   `db:"price" json:"price"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func (ic *ItemCreate) Validate() error {
	if len(ic.Name) < 3 || len(ic.Name) > 128 {
		return NameLengthIncorrect
	} else if len(ic.Name) > 1024 {
		return DescriptionLengthIncorrect
	}

	return nil
}

func (ic *ItemCreate) PrepareCreate() error {
	ic.Name = strings.TrimSpace(ic.Name)
	ic.Description = strings.TrimSpace(ic.Description)

	err := ic.Validate()
	if err != nil {
		return err
	}

	return nil
}
