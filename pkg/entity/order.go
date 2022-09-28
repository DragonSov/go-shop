package entity

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type OrderCreate struct {
	Items map[uuid.UUID]int `db:"items" json:"items"`
}

type OrderItems map[uuid.UUID]int

type Order struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Items     OrderItems `db:"items" json:"items"`
	Price     float32    `db:"price" json:"price"`
	OrderedBy uuid.UUID  `db:"ordered_by" json:"ordered_by"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

func (oi *OrderItems) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &oi)
		if err != nil {
			return err
		}
	}

	return nil
}

func (oi *OrderItems) Value() (driver.Value, error) {
	return json.Marshal(&oi)
}
