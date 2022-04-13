package models

import (
	"context"
	"fmt"
	"github.com/phamtrunghieu/tinder-clone-backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type User struct {
	Uuid        string    `json:"uuid,omitempty" bson:"uuid"`
	Id          string    `json:"id" bson:"id"`
	LastName    string    `json:"last_name,omitempty" bson:"last_name"`
	FirstName   string    `json:"first_name" bson:"first_name"`
	Gender      string    `json:"gender" bson:"gender"`
	Picture     string    `json:"picture" bson:"picture"`
	DateOfBirth string    `json:"date_of_birth" bson:"date_of_birth"`
	Email       string    `json:"email" bson:"email"`
	Phone       string    `json:"phone" bson:"phone"`
	Matches     []string  `json:"matches" bson:"matches"`
	CreatedAt   time.Time `json:"-" bson:"created_at"`
	UpdatedAt   time.Time `json:"-" bson:"updated_at"`
}

func (u *User) Model() *mongo.Collection {
	store := database.GetInstance().Collection("user")

	return store
}

func (u *User) Insert() error {
	_, err := u.Model().InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Find(ctx context.Context, condition map[string]interface{}) ([]User, error) {
	user := u.Model()
	cur, err := user.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var result []User
	for cur.Next(context.Background()) {
		var item User
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

func (u *User) FindOne(ctx context.Context, condition map[string]interface{}) (*User, error) {
	user := u.Model()
	var result *User
	err := user.FindOne(ctx, condition).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *User) FindRandomUser(ctx context.Context) (*User, error) {
	user := u.Model()
	pipeline := []bson.M{
		{
			"$sample": bson.M{"size": 1},
		},
	}
	cur, err := user.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var result *User
	for cur.Next(ctx) {
		var item User
		err := cur.Decode(&item)
		if err != nil {
			fmt.Println("[ERROR][Decode]", err)
			fmt.Println("[ERROR][Item]", item)
			continue
		}
		result = &item
	}

	return result, nil
}

func (u *User) UpdatePushDataToArray(conditions map[string]interface{}, data map[string]interface{}) (int64, error) {
	coll := u.Model()

	updateStr := make(map[string]interface{})
	updateStr["$addToSet"] = data
	resp, err := coll.UpdateOne(context.TODO(), conditions, updateStr)
	if err != nil {
		log.Println("Can not update StoreGroup", err)
		return 0, err
	}

	return resp.ModifiedCount, nil
}
