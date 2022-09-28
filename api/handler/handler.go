package handler

import (
	"errors"
	"lessons/pkg/entity"
	"lessons/pkg/item"
	"lessons/pkg/order"
	"lessons/pkg/user"
	"log"
	"strings"
)

var (
	AuthTokenIncorrect       = errors.New("неверный HTTP-заголовок авторизации")
	BodyParsingError         = errors.New("произошла ошибка при обработке тела запроса")
	NotEnoughPermissions     = errors.New("у вас недостаточно прав для выполнения этого действия")
	UUIDIncorrect            = errors.New("укажите корректный UUID")
	LoginOrPasswordIncorrect = errors.New("неверный логин или пароль")
	QueryParsingError        = errors.New("произошла ошибка при обработке параметров запроса")
)

func exceptionHandler(err error) (error, int) {
	switch err {
	case AuthTokenIncorrect, NotEnoughPermissions, LoginOrPasswordIncorrect, user.JWTError:
		return err, 403
	case item.NotFound, user.NotFound, order.NotFound, order.OneOfItemNotFound:
		return err, 404
	case user.AlreadyExists:
		return err, 409
	case BodyParsingError, UUIDIncorrect, QueryParsingError, entity.LimitLengthIncorrect, entity.OffsetLengthIncorrect,
		entity.DescriptionLengthIncorrect, entity.NameLengthIncorrect, entity.LoginContentIncorrect,
		entity.LoginLengthIncorrect, entity.PasswordLengthIncorrect:
		return err, 422
	default:
		log.Println(err)
		return errors.New("произошла непредвиденная ошикба"), 500
	}
}

func getCurrentUser(userService user.Service, authHeader string) (*entity.User, error) {
	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		return nil, AuthTokenIncorrect
	}

	userData, err := userService.GetCurrentUser(authToken[1])
	if err != nil {
		return nil, err
	}

	return userData, nil
}
