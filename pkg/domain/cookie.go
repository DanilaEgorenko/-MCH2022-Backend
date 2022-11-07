package domain

type CookieStorage interface {
	Add(key string, userId int) error
	Delete(value string) error
}

type CookieUseCase interface {
	Add(key string, userId int) error
	Delete(value string) error
}
