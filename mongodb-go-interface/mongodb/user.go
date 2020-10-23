package mongodb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Name       string             `json:"name,omitempty"`
	Email      string             `json:"email,omitempty"`
	Password   string             `json:"password,omitempty"`
	SettingsID primitive.ObjectID `bson:"_settings_id" json:"settings_id,omitempty"`
}

// AddUser creates a new user entry in the DB.
// Whitespaces in the email are automatically deleted
// Email is a unique attribute, so the function checks for existing email, before adding a new entry
func AddUser(name string, email string, passwd string) (*mongo.InsertOneResult, error) {
	email = strings.ReplaceAll(email, " ", "")

	client, err := Connect()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd), 16)
	if err != nil {
		return nil, err
	}

	settingsID, err := AddSettings()
	if err != nil {
		return nil, err
	}

	user := User{
		Name:       name,
		Email:      email,
		Password:   string(hashedPassword),
		SettingsID: settingsID.InsertedID.(primitive.ObjectID),
	}

	existingUser, err := GetUserByEmail(email)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			return nil, err
		}

		if existingUser.Email == email {
			return nil, fmt.Errorf("User with %s email already exists", email)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("user_database").Collection("users")
	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect mongo client: %s\n", errors.WithStack(err))
	}

	return insertResult, nil
}

// GetUserByEmail returns the user defined by the email and password.
// Whitespaces in the email are automatically removed.
func GetUserByEmail(email string) (User, error) {
	email = strings.ReplaceAll(email, " ", "")

	client, err := Connect()
	if err != nil {
		return User{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("user_database").Collection("users")
	result := collection.FindOne(ctx, bson.M{"email": email})
	user := User{}
	if err := result.Decode(&user); err != nil {
		return User{}, err
	}

	if err = client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect mongo client: %s\n", errors.WithStack(err))
	}

	return user, nil
}

func CheckEmailAndPassword(email string, password string) error {
	user, err := GetUserByEmail(email)
	if err == mongo.ErrNoDocuments {
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

// DeleteUser removes the user defined by the email address.
// Whitespaces in the email are automatically removed.
func deleteUserEntry(email string) (*mongo.DeleteResult, error) {
	email = strings.ReplaceAll(email, " ", "")

	client, err := Connect()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("user_database").Collection("users")
	result, err := collection.DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}

	if err = client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect mongo client: %s\n", errors.WithStack(err))
	}

	return result, nil
}

func DeleteUser(email string) error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return err
	}

	_, err = deleteUserEntry(email)
	if err != nil {
		return err
	}

	result, err := DeleteSettings(&user.SettingsID)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		log.Println("There is no settings associated with the user ", email)
	}
	return nil
}
