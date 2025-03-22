package watchParties

import (
	"crypto/rand"
	"wat/internal/database"
	"wat/internal/models"
)

func CreateWatchParty(db database.DatabaseController, watchParty *models.WatchParty) (string, error) {
	watchPartyId := rand.Text()
	err := db.CreateWatchParty(watchPartyId, watchParty)
	if err != nil {
		return "", err
	}
	return watchPartyId, nil
}

func GetAllWatchParties(db database.DatabaseController) (map[string]*models.WatchParty, error) {
	watchParties, err := db.GetAllWatchParties()
	if err != nil {
		return nil, err
	}
	return watchParties, nil
}

func DeleteWatchParty(db database.DatabaseController, watchPartyId string) error {
	err := db.DeleteWatchParty(watchPartyId)
	if err != nil {
		return err
	}
	return nil
}
