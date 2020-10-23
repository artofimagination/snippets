package mysqldb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User defines the user structures. Each user must have an associated settings entry.
type User struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password,omitempty"`
	SettingsID uuid.UUID `json:"user_settings_id,omitempty"`
}

// GetUserByEmail returns the user defined by the email and password.
func GetUserByEmail(email string) (User, error) {
	email = strings.ReplaceAll(email, " ", "")

	var user User
	queryString := "select BIN_TO_UUID(id), name, email, password, BIN_TO_UUID(user_settings_id) from users where email = ?"
	db, err := ConnectSystem()
	if err != nil {
		return user, err
	}
	defer db.Close()

	query, err := db.Query(queryString, email)

	if err != nil {
		return user, err
	}
	defer query.Close()

	query.Next()
	if err := query.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.SettingsID); err != nil {
		return user, err
	}

	return user, nil
}

func UserExists(username string) (bool, error) {
	var user User
	db, err := ConnectSystem()
	if err != nil {
		return false, err
	}
	defer db.Close()

	queryString := "SELECT name FROM users WHERE name = ?"
	queryUser := db.QueryRow(queryString, username)
	err = queryUser.Scan(&user.Name)
	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		return false, err
	default:
		return true, err
	}

	return false, nil
}

func EmailExists(email string) (bool, error) {
	email = strings.ReplaceAll(email, " ", "")

	var user User
	queryString := "select email from users where email = ?"
	db, err := ConnectSystem()
	if err != nil {
		return false, err
	}
	defer db.Close()

	queryEmail := db.QueryRow(queryString, email)

	err = queryEmail.Scan(&user.Email)
	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		return false, err
	default:
		return true, err
	}

	return false, nil
}

// CheckPassword compares the password entered by the user with the stored password.
func CheckEmailAndPassword(email string, password string) error {
	email = strings.ReplaceAll(email, " ", "")

	user, err := GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return fmt.Errorf("Incorrect email or password")
	}

	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return fmt.Errorf("Incorrect email or password")
	}
	return nil
}

// AddUser creates a new user entry in the DB.
// Whitespaces in the email are automatically deleted
// Email is a unique attribute, so the function checks for existing email, before adding a new entry
func AddUser(name string, email string, passwd string) error {
	email = strings.ReplaceAll(email, " ", "")

	queryString := "INSERT INTO users (id, name, email, password, user_settings_id) VALUES (UUID_TO_BIN(UUID()), ?, ?, ?, UUID_TO_BIN(?))"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}

	defer db.Close()

	settingsID, err := AddSettings()
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd), 16)
	if err != nil {
		if err := DeleteSettings(settingsID); err != nil {
			return errors.Wrap(errors.WithStack(err), "Failed to revert settings creation")
		}
		return err
	}

	query, err := db.Query(queryString, name, email, hashedPassword, &settingsID)
	if err != nil {
		if err := DeleteSettings(settingsID); err != nil {
			return errors.Wrap(errors.WithStack(err), "Failed to revert settings creation")
		}
		return err
	}

	defer query.Close()
	return nil
}

func deleteUserEntry(email string) error {
	query := "DELETE FROM users WHERE email=?"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query, email)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(email string) error {
	email = strings.ReplaceAll(email, " ", "")
	user, err := GetUserByEmail(email)
	if err != nil {
		return err
	}

	if err := deleteUserEntry(email); err != nil {
		return err
	}

	if err := DeleteSettings(&user.SettingsID); err != nil {
		return err
	}

	return nil
}
