package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryDB struct {
	collection *mongo.Collection
}

func NewUserRepositoryDB(client *mongo.Client, dbName string, collName string) UserRepository {
	collection := client.Database(dbName).Collection(collName)
	return &userRepositoryDB{collection: collection}
}

func (r *userRepositoryDB) InsertUser(user *User) error {
	ctx := context.Background()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r userRepositoryDB) GetUser(gId string) (*User, error) {
	ctx := context.Background()
	var user User
	filter := bson.M{"googleId": gId}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
