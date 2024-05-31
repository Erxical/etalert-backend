package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type bedtimeRepositoryDB struct {
	collection *mongo.Collection
}

func NewBedtimeRepositoryDB(client *mongo.Client, dbName string, collName string) BedtimeRepository {
	collection := client.Database(dbName).Collection(collName)
	return &bedtimeRepositoryDB{collection: collection}
}

func (r *bedtimeRepositoryDB) InsertBedtime(bedtime *Bedtime) error {
	ctx := context.Background()
	_, err := r.collection.InsertOne(ctx, bedtime)
	return err
}

func (r bedtimeRepositoryDB) GetBedtimeInfo(gId string) (*Bedtime, error) {
	ctx := context.Background()
	var bedtime Bedtime
	filter := bson.M{"googleId": gId}
	err := r.collection.FindOne(ctx, filter).Decode(&bedtime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &bedtime, nil
}

func (r *bedtimeRepositoryDB) UpdateBedtime(gId string, bedtime *Bedtime) error {
	ctx := context.Background()
	filter := bson.M{"googleId": gId}
	update := bson.M{"$set": bedtime}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
