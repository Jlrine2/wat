package auth

import (
	"crypto/rand"
	"wat/internal/database"
	"wat/internal/models"
)

func CreateSession(accessTokenDetails *models.AccessTokenDetails, db database.DatabaseController) (string, error) {
	sessionId := rand.Text()
	err := db.SaveAuthSession(sessionId, accessTokenDetails)
	if err != nil {
		return "", err
	}
	return sessionId, nil
}

func GetSession(sessionId string, db database.DatabaseController) (*models.AccessTokenDetails, error) {
	value, err := db.GetAuthSession(sessionId)
	if err != nil {
		return nil, err
	}
	return value, nil
}
