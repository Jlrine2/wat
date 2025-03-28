package database

import (
	"fmt"
	"wat/internal/models"
)

type MemoryDatabase struct {
	sessions     map[string]*models.AccessTokenDetails
	watchParties map[string]*models.WatchParty
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		sessions:     make(map[string]*models.AccessTokenDetails),
		watchParties: make(map[string]*models.WatchParty),
	}
}

func (db *MemoryDatabase) SaveAuthSession(key string, value *models.AccessTokenDetails) error {
	db.sessions[key] = value
	return nil
}

func (db *MemoryDatabase) GetAuthSession(key string) (*models.AccessTokenDetails, error) {
	result, ok := db.sessions[key]
	if !ok {
		return nil, fmt.Errorf("auth session not found")
	}

	return result, nil
}

func (db *MemoryDatabase) CreateWatchParty(key string, value *models.WatchParty) error {
	db.watchParties[key] = value
	return nil
}

func (db *MemoryDatabase) GetAllWatchParties() (map[string]*models.WatchParty, error) {
	return db.watchParties, nil
}

func (db *MemoryDatabase) DeleteWatchParty(key string) error {
	delete(db.watchParties, key)
	return nil
}
