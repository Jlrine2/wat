package auth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"wat/internal/database"
)

func CreateSession(accessTokenDetails *DiscordAccessTokenDetails, db database.DatabaseController) (string, error) {
	sessionId := rand.Text()
	accessTokenBytes, err := json.Marshal(accessTokenDetails)
	if err != nil {
		return "", fmt.Errorf("Unable to encode access token as bytes: %s", err.Error())
	}
	err = db.SaveAuthSession(sessionId, string(accessTokenBytes))
	if err != nil {
		return "", fmt.Errorf("unable to save auth session in DB: %s", err.Error())
	}
	return sessionId, nil
}

func GetSession(sessionId string, db database.DatabaseController) (*DiscordAccessTokenDetails, error) {
	value, err := db.GetAuthSession(sessionId)
	if err != nil {
		return nil, err
	}
	var accessToken DiscordAccessTokenDetails
	err = json.Unmarshal([]byte(value), &accessToken)
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}
