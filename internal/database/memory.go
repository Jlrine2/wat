package database

import (
	"fmt"
	"wat/internal/models"
)

type MemoryDatabase struct {
	sessions map[string]any
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		sessions: map[string]any{},
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

	return result.(*models.AccessTokenDetails), nil
}
