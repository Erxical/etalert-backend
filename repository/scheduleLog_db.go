package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type scheduleLogRepositoryDB struct {
	collection *mongo.Collection
}

func NewScheduleLogRepositoryDB(client *mongo.Client, dbName string, collName string) ScheduleLogRepository {
	collection := client.Database(dbName).Collection(collName)
	return &scheduleLogRepositoryDB{collection: collection}
}

func (s *scheduleLogRepositoryDB) InsertScheduleLog(scheduleLog *ScheduleLog) error {
	ctx := context.Background()
	_, err := s.collection.InsertOne(ctx, scheduleLog)
	return err
}