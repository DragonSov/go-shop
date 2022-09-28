package item

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"lessons/pkg/entity"
)

type Repository interface {
	Create(item *entity.ItemCreate) (*entity.Item, error)
	Get(limit int, offset int) (*[]entity.Item, error)
	GetByID(id uuid.UUID) (*entity.Item, error)
}
type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(item *entity.ItemCreate) (*entity.Item, error) {
	var newItem entity.Item
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	err = tx.Get(&newItem, "INSERT INTO items (name, description, price) VALUES ($1, $2, $3) RETURNING *", item.Name, item.Description, item.Price)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &newItem, nil
}

func (r *repository) Get(limit int, offset int) (*[]entity.Item, error) {
	var selectedItems []entity.Item
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	err = tx.Select(&selectedItems, "SELECT * FROM items LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}

	return &selectedItems, nil
}

func (r *repository) GetByID(id uuid.UUID) (*entity.Item, error) {
	var selectedItem entity.Item
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	err = tx.Get(&selectedItem, "SELECT * FROM items WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &selectedItem, nil
}
