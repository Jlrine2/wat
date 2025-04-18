package database

import "wat/internal/models"

type DatabaseController interface {
	SaveAuthSession(key string, value *models.Session) error
	GetAuthSession(key string) (*models.Session, error)

	CreateWatchParty(key string, watchParty *models.WatchParty) error
	GetAllWatchParties() (map[string]*models.WatchParty, error)
	DeleteWatchParty(key string) error
}
