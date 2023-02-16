package repository

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ScopeActivation = "activation"
	ScopeAuthentication = "authentication"
	ScopePasswordReset = "password-reset"
)

type Token struct {
	Scope string `-`
	UserID string `-`
	Expiry time.Time `json:"expiry"`
	PlainText string `json:"token"`
	Hash []byte `-`
}

type TokenUser struct {
	Scope string `-`
	UserID string `-`
	Expiry time.Time `json:"expiry"`
	PlainText string `json:"token"`
	Hash []byte `-`
	User []User
}

func generateToken() {}

type TokenModel struct {db *mongo.Collection}

func (m *TokenModel) Insert(token *Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	doc, err := bson.Marshal(token)
	if err != nil {
		return err
	}

	_, err = m.db.InsertOne(ctx,doc)
	if err != nil {
		return err
	}

	return nil
}

func (m *TokenModel) GetUser(tokenPlainText, scope string) (*TokenUser, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	filter := bson.D{
		{"hash", tokenHash}, 
		{"scope", scope}, 
		{"expire", bson.D{{"$gte", time.Now()}}},
	}

	sortStage := bson.D{{"$sort", bson.D{{"expiry", 1}}}}

	limitStage := bson.D{{"$limit", 1}}

	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "users"}, 
			{"localField", "userID"},
			{"foreignField", "_id"},
			{"as", "user"},
		}},
	}

	// unwindStage := bson.M{
	// 	"$unwind": bson.M{
	// 		"path": "$user",
	// 		"preserveNullAndEmptyArrays": false,
	// 	},
	// }

	pipeline := mongo.Pipeline{filter, sortStage, limitStage, lookupStage}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.db.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err 
	}

	var token []TokenUser

	if err = cursor.All(ctx, &token); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	return &token[0], nil
}

func (m *TokenModel) Get(tokenPlainText, scope string) (*Token, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	filter := bson.M{
		"hash": tokenHash, 
		"scope": scope, 
		"expire": bson.M{"$gte": time.Now()},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var token Token
	err := m.db.FindOne(ctx, filter).Decode(&token)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	return &token, nil
}


func (m *TokenModel) DeleteAllForUser(id, scope string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{
		"userID": id,
		"scope": scope,
	}

	_, err := m.db.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
} 