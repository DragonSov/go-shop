package order

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"lessons/pkg/entity"
)

type Repository interface {
	Create(order *entity.Order) (*entity.Order, error)
	Get(limit int, offset int, orderedBy uuid.UUID) (*[]entity.Order, error)
	GetByID(id uuid.UUID) (*entity.Order, error)
}
type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(order *entity.Order) (*entity.Order, error) {
	var newOrder entity.Order
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return nil, err
	}

	err = tx.Get(&newOrder, "INSERT INTO orders (items, price, ordered_by) VALUES ($1, $2, $3) RETURNING *", itemsJSON, order.Price, order.OrderedBy)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &newOrder, nil
}

func (r *repository) Get(limit int, offset int, orderedBy uuid.UUID) (*[]entity.Order, error) {
	var selectedOrders []entity.Order
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	err = tx.Select(&selectedOrders, "SELECT * FROM orders WHERE ordered_by = $1 LIMIT $2 OFFSET $3", orderedBy, limit, offset)
	if err != nil {
		return nil, err
	}

	return &selectedOrders, nil
}

func (r *repository) GetByID(id uuid.UUID) (*entity.Order, error) {
	var selectedOrder entity.Order
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	err = tx.Get(&selectedOrder, "SELECT * FROM orders WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &selectedOrder, nil
}
