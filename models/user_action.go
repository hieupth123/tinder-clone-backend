package models

import (
	"context"
	"fmt"
	"github.com/phamtrunghieu/tinder-clone-backend/database"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserAction struct {
	Uuid      string    `json:"uuid,omitempty" bson:"uuid"`
	UserUuid  string    `json:"user_uuid" bson:"user_uuid"`
	Type      string    `json:"type" bson:"type"`
	GuestUuid string    `json:"guest_uuid" bson:"guest_uuid"`
	CreatedAt time.Time `json:"-" bson:"created_at"`
	UpdatedAt time.Time `json:"-" bson:"updated_at"`
}

func (u *UserAction) Model() *mongo.Collection {
	store := database.GetInstance().Collection("user_action")

	return store
}

func (u *UserAction) FindOne(ctx context.Context, condition map[string]interface{}) (*UserAction, error) {
	user := u.Model()
	var result *UserAction
	err := user.FindOne(ctx, condition).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *UserAction) Insert() error {
	_, err := u.Model().InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserAction) Find(ctx context.Context, condition map[string]interface{}) ([]UserAction, error) {
	userAction := u.Model()
	cur, err := userAction.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var result []UserAction
	for cur.Next(context.Background()) {
		var item UserAction
		err := cur.Decode(&item)

		if err != nil {
			fmt.Println("[ERROR][Decode]", err)
			fmt.Println("[ERROR][Item]", item)
			continue
		}
		result = append(result, item)
	}

	return result, nil
}
