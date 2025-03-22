package database

import "fmt"

type MemoryDatabase struct {
	sessions map[string]string
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		sessions: map[string]string{},
	}
}

func (db *MemoryDatabase) SaveAuthSession(key string, value string) error {
	db.sessions[key] = value
	return nil
}

func (db *MemoryDatabase) GetAuthSession(key string) (string, error) {
	result, ok := db.sessions[key]
	if !ok {
		return "", fmt.Errorf("auth session not found")
	}
	return result, nil
}
