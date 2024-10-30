package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type feedbackRepositoryDB struct {
	collection *mongo.Collection
}

func NewFeedbackRepositoryDB(client *mongo.Client, dbName string, collName string) FeedbackRepository {
	collection := client.Database(dbName).Collection(collName)
	return &feedbackRepositoryDB{collection: collection}
}

func (f *feedbackRepositoryDB) InsertFeedback(feedback *Feedback) error {
	ctx := context.Background()
	_, err := f.collection.InsertOne(ctx, feedback)
	return err
}