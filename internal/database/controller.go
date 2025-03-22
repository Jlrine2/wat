package database

type DatabaseController interface {
	SaveAuthSession(key string, value string) error
	GetAuthSession(key string) (string, error)
}
