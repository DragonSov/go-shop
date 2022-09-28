package user

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"lessons/pkg/entity"
)

type Repository interface {
	Create(user *entity.User) (*entity.User, error)
	GetByLogin(login string) (*entity.User, error)
	GetByID(id uuid.UUID) (*entity.User, error)
	UpdatePasswordByID(id uuid.UUID, password string) (*entity.User, error)
}
type repository struct {
	Database *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		Database: db,
	}
}

func (r *repository) Create(user *entity.User) (*entity.User, error) {
	var newUser entity.User
	tx, err := r.Database.Beginx()
	if err != nil {
		return nil, err
	}

	rows := tx.QueryRowx("INSERT INTO users (login, password) VALUES ($1, $2) RETURNING *", user.Login, user.Password)
	if err = rows.StructScan(&newUser); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return &newUser, nil
}

func (r *repository) GetByLogin(login string) (*entity.User, error) {
	var user entity.User
	tx, err := r.Database.Beginx()
	if err != nil {
		return nil, err
	}

	err = tx.Get(&user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	tx, err := r.Database.Beginx()
	if err != nil {
		return nil, err
	}

	err = tx.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) UpdatePasswordByID(id uuid.UUID, password string) (*entity.User, error) {
	var user entity.User
	tx, err := r.Database.Beginx()
	if err != nil {
		return nil, err
	}

	err = tx.Get(&user, "UPDATE users SET password = $1 WHERE id = $2 RETURNING *", password, id)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}
