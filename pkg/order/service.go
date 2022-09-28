package order

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"lessons/pkg/entity"
	"lessons/pkg/item"
)

var (
	NotFound          = errors.New("данный заказ не найден")
	OneOfItemNotFound = errors.New("один из предметов заказа не найден")
)

type Service interface {
	Create(order *entity.Order) (*entity.Order, error)
	Get(limit int, offset int, orderedBy uuid.UUID) (*[]entity.Order, error)
	GetByID(id uuid.UUID) (*entity.Order, error)
}

type service struct {
	itemRepository  item.Repository
	orderRepository Repository
}

func NewService(itemRepo item.Repository, orderRepo Repository) Service {
	return &service{
		itemRepository:  itemRepo,
		orderRepository: orderRepo,
	}
}

func (s *service) Create(order *entity.Order) (*entity.Order, error) {
	for itemID, count := range order.Items {
		itemData, err := s.itemRepository.GetByID(itemID)

		if err == sql.ErrNoRows || itemData == nil {
			return nil, OneOfItemNotFound
		} else if err != nil {
			return nil, err
		}

		order.Price += itemData.Price * float32(count)
	}

	return s.orderRepository.Create(order)
}

func (s *service) Get(limit int, offset int, orderedBy uuid.UUID) (*[]entity.Order, error) {
	orders, err := s.orderRepository.Get(limit, offset, orderedBy)
	if err == sql.ErrNoRows {
		return nil, NotFound
	} else if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *service) GetByID(id uuid.UUID) (*entity.Order, error) {
	order, err := s.orderRepository.GetByID(id)
	if err == sql.ErrNoRows {
		return nil, NotFound
	} else if err != nil {
		return nil, err
	}

	return order, nil
}
