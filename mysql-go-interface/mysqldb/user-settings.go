package mysqldb

import (
	"log"

	"github.com/google/uuid"
)

type UserSettings struct {
	ID                uuid.UUID `json:"id,omitempty"`
	Enable2StepsVerif bool      `json:"two_steps_verif,omitempty"`
}

func AddSettings() (*uuid.UUID, error) {
	queryString := "INSERT INTO user_settings (id, two_steps_verif) VALUES (UUID_TO_BIN(?), ?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	newID, err := uuid.NewUUID()
	log.Println(newID)
	if err != nil {
		return nil, err
	}

	query, err := db.Query(queryString, newID, false)
	if err != nil {
		return nil, err
	}

	defer query.Close()
	return &newID, nil
}

func GetSettings(settingsID *uuid.UUID) (*UserSettings, error) {
	settings := UserSettings{}
	queryString := "SELECT two_steps_verif FROM user_settings WHERE id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := db.QueryRow(queryString, *settingsID)

	if err != nil {
		return nil, err
	}

	if err := query.Scan(&settings.Enable2StepsVerif); err != nil {
		return nil, err
	}

	return &settings, nil
}

func DeleteSettings(settingsID *uuid.UUID) error {
	query := "DELETE FROM user_settings WHERE id=UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query, *settingsID)
	if err != nil {
		return err
	}
	return nil
}
