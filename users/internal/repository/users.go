package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersModel struct {db *mongo.Collection}

type User struct {
	ID string `json:"_id"`
	CreatedAt time.Time `json:"created_at"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	Password password `json:"password"`
	Username string `json:"username"`
}

type password struct {
	plaintext *string
	hash []byte
}

func (m UsersModel) Insert(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.db.InsertOne(ctx, user)
	return err
}

func (m UsersModel) GetByFilter(data map[string]any) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	user := User{}
	filter := bson.M(data)

	err := m.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}