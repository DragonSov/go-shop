package item

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"lessons/pkg/entity"
)

var (
	NotFound = errors.New("данный предмет не найден")
)

type Service interface {
	Create(item *entity.ItemCreate) (*entity.Item, error)
	Get(limit int, offset int) (*[]entity.Item, error)
	GetByID(id uuid.UUID) (*entity.Item, error)
}

type service struct {
	itemRepository Repository
}

func NewService(itemRepo Repository) Service {
	return &service{
		itemRepository: itemRepo,
	}
}

func (s *service) Create(item *entity.ItemCreate) (*entity.Item, error) {
	return s.itemRepository.Create(item)
}

func (s *service) Get(limit int, offset int) (*[]entity.Item, error) {
	return s.itemRepository.Get(limit, offset)
}

func (s *service) GetByID(id uuid.UUID) (*entity.Item, error) {
	item, err := s.itemRepository.GetByID(id)
	if err == sql.ErrNoRows {
		return nil, NotFound
	} else if err != nil {
		return nil, err
	}

	return item, nil
}
