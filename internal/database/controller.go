package database

import "wat/internal/models"

type DatabaseController interface {
	SaveAuthSession(key string, value *models.AccessTokenDetails) error
	GetAuthSession(key string) (*models.AccessTokenDetails, error)
}
