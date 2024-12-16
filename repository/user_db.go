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

func (r userRepositoryDB) GetUserInfo(gId string) (*User, error) {
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

func (r userRepositoryDB) UpdateUser(gId string, user *User) error {
	ctx := context.Background()
	filter := bson.M{"googleId": gId}
	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"image": user.Image,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r userRepositoryDB) GetAllUsersId() ([]string, error) {
	ctx := context.Background()
	var users []string
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user.GoogleId)
	}
	return users, nil
}