package user

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"lessons/pkg/common/config"
	"lessons/pkg/common/utils"
	"lessons/pkg/entity"
)

type Service interface {
	Create(user *entity.User) (*entity.User, error)
	GetCurrentUser(token string) (*entity.User, error)
	GetByID(userID uuid.UUID) (*entity.User, error)
	GetByLogin(login string) (*entity.User, error)
	UpdatePasswordByID(userID uuid.UUID, password string) (*entity.User, error)
}
type service struct {
	userRepository Repository
	cfg            *config.Config
}

var (
	AlreadyExists = errors.New("пользователь с данным логином уже существует")
	NotFound      = errors.New("пользователь не найден")
	UUIDIncorrect = errors.New("неккоректный UUID")
	JWTError      = errors.New("неккоректный токен авторизации")
)

func NewService(r Repository, cfg *config.Config) Service {
	return &service{
		userRepository: r,
		cfg:            cfg,
	}
}

func (s *service) Create(user *entity.User) (*entity.User, error) {
	selectedUser, err := s.GetByLogin(user.Login)
	if selectedUser != nil {
		return nil, AlreadyExists
	} else if err != nil && err != NotFound {
		return nil, err
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, err
	}
	return s.userRepository.Create(user)
}

func (s *service) GetCurrentUser(token string) (*entity.User, error) {
	claims, err := utils.DecodeJWTToken(token, s.cfg)
	if err != nil {
		return nil, JWTError
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, UUIDIncorrect
	}

	selectedUser, err := s.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return selectedUser, nil
}

func (s *service) GetByLogin(login string) (*entity.User, error) {
	user, err := s.userRepository.GetByLogin(login)
	if err == sql.ErrNoRows {
		return nil, NotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetByID(userID uuid.UUID) (*entity.User, error) {
	user, err := s.userRepository.GetByID(userID)
	if err == sql.ErrNoRows {
		return nil, NotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdatePasswordByID(userID uuid.UUID, password string) (*entity.User, error) {
	selectedUser, err := s.GetByID(userID)
	if selectedUser == nil {
		return nil, NotFound
	} else if err != nil {
		return nil, err
	}

	selectedUser.Password = password
	if err = selectedUser.PrepareUpdate(); err != nil {
		return nil, err
	}

	return s.userRepository.UpdatePasswordByID(selectedUser.ID, selectedUser.Password)
}
