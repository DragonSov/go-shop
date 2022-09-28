package entity

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserRole string

const (
	CustomerRole UserRole = "customer"
	AdminRole             = "admin"
)

var (
	LoginContentIncorrect   = errors.New("логин должен состоять только из цифр или символов английского алфавита")
	LoginLengthIncorrect    = errors.New("логин должен быть длиной от 3 до 64 символов")
	PasswordLengthIncorrect = errors.New("пароль должен быть длиной от 6 до 128 символов")
)

type UserCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Login     string    `db:"login" json:"login"`
	Password  string    `db:"password" json:"password"`
	Role      UserRole  `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

func checkLoginAlphabet(login string) bool {
	for _, char := range login {
		checkEnglishLetter := char >= 'a' && char <= 'z'
		checkNumber := char >= '0' && char <= '9'
		if !checkEnglishLetter && !checkNumber {
			return false
		}
	}

	return true
}

func (uc *UserCredentials) Validate() error {
	if loginAlphabetValid := checkLoginAlphabet(uc.Login); !loginAlphabetValid {
		return LoginContentIncorrect
	} else if len(uc.Login) < 3 || len(uc.Login) > 64 {
		return LoginLengthIncorrect
	}

	if len(uc.Password) < 6 || len(uc.Password) > 128 {
		return PasswordLengthIncorrect
	}

	return nil
}

func (uc *UserCredentials) ToLowerCase() error {
	uc.Login = strings.ToLower(uc.Login)
	if err := uc.Validate(); err != nil {
		return err
	}

	return nil
}

func (u *User) Validate() error {
	if loginAlphabetValid := checkLoginAlphabet(u.Login); !loginAlphabetValid {
		return LoginContentIncorrect
	} else if len(u.Login) < 3 || len(u.Login) > 64 {
		return LoginLengthIncorrect
	}

	if len(u.Password) < 6 || len(u.Password) > 128 {
		return PasswordLengthIncorrect
	}

	return nil
}

func (u *User) PrepareCreate() error {
	u.Login = strings.ToLower(u.Login)
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

func (u *User) PrepareUpdate() error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
