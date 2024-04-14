package storage

import "errors"

var (
	//TODO
	ErrContentNotFound = errors.New("Баннер для не найден")
	ErrInternalServer  = errors.New("Внутренняя ошибка сервера")
	ErrIncorrectData   = errors.New("Некорректные данные")
	ErrBannerNotFound  = errors.New("Баннер не найден")
	ErrUserNotFound    = errors.New("Пользователь не авторизован")
)
