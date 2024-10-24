package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type routineLogRepositoryDB struct {
	collection *mongo.Collection
}

func NewRoutineLogRepositoryDB(client *mongo.Client, dbName string, collName string) RoutineLogRepository {
	collection := client.Database(dbName).Collection(collName)
	return &routineLogRepositoryDB{collection: collection}
}

func (r *routineLogRepositoryDB) InsertRoutineLog(routineLog *RoutineLog) error {
	ctx := context.Background()
	_, err := r.collection.InsertOne(ctx, routineLog)
	return err
}

func (r *routineLogRepositoryDB) GetRoutineLogs(googleId string, date string) ([]*RoutineLog, error) {
	ctx := context.Background()
	var routineLogs []*RoutineLog

	filter := bson.M{"googleId": googleId}
	if date != "" {
		filter["date"] = date
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var routineLog RoutineLog
		if err := cursor.Decode(&routineLog); err != nil {
			return nil, err
		}
		routineLogs = append(routineLogs, &routineLog)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return routineLogs, nil
}

func (r *routineLogRepositoryDB) DeleteRoutineLog(id string) error {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ID: %v", err)
	}

	filter := bson.M{"_id": objectId}
	_, err = r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete routine log: %v", err)
	}

	return nil
}
