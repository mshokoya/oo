package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)



type User struct {
	ID string `json:"_id"`
	CreatedAt time.Time `json:"created_at"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	Password Password `json:"password"`
	Username string `json:"username"`
	Email string `json:"email"`
}

var UserList = []string{"id", "created_at", "firstname", "lastname", "password", "username", "email"}

type UsersModel struct {db *mongo.Collection}

func (m UsersModel) Insert(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.db.InsertOne(ctx, user)
	return err
}

func (m UsersModel) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	user := User{}

	filter := bson.M{"email": email}

	err := m.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m UsersModel) Update(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.ID}

	update, err := bson.Marshal(user)
	if err != nil {
		return err
	}

	_, err = m.db.UpdateOne(ctx, filter, bson.M{"$set": update});

	if err != nil {
		return err
	}

	return nil
}

type Password struct {
	Plaintext *string
	Hash []byte
}

func (p *Password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &password
	p.Hash = hash

	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}