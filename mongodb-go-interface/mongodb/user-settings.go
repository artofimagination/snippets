package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserSettings struct {
	Enable2StepsVerif bool `json:"2steps_on,omitempty"`
}

func AddSettings() (*mongo.InsertOneResult, error) {
	client, err := Connect()
	if err != nil {
		return nil, err
	}

	userSettings := UserSettings{
		Enable2StepsVerif: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("user_database").Collection("user_settings")
	insertResult, err := collection.InsertOne(ctx, userSettings)
	if err != nil {
		return nil, err
	}

	if err = client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect mongo client: %s\n", errors.WithStack(err))
	}
	return insertResult, nil
}

func GetSettings(objID *primitive.ObjectID) (UserSettings, error) {
	if objID == nil {
		return UserSettings{}, fmt.Errorf("Invalid settings ID")
	}

	client, err := Connect()
	if err != nil {
		return UserSettings{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("user_database").Collection("user_settings")
	result := collection.FindOne(ctx, bson.M{"_id": objID})
	settings := UserSettings{}
	if err := result.Decode(&settings); err != nil {
		return UserSettings{}, err
	}

	if err = client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect mongo client: %s\n", errors.WithStack(err))
	}

	return settings, nil
}

func DeleteSettings(objID *primitive.ObjectID) (*mongo.DeleteResult, error) {
	if objID == nil {
		return nil, fmt.Errorf("Invalid settings ID")
	}

	client, err := Connect()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("user_database").Collection("user_settings")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return nil, err
	}

	if err = client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect mongo client: %s\n", errors.WithStack(err))
	}
	return result, nil
}
