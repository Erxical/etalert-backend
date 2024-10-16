package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type scheduleLogRepositoryDB struct {
	collection *mongo.Collection
}

func NewScheduleLogRepositoryDB(client *mongo.Client, dbName string, collName string) ScheduleLogRepository {
	collection := client.Database(dbName).Collection(collName)
	return &scheduleLogRepositoryDB{collection: collection}
}

func (s *scheduleLogRepositoryDB) GetUpcomingSchedules() ([]int, error) {
	ctx := context.Background()
	var groupIds []int

	now := time.Now().UTC().Add(7 * time.Hour)
	currentTime := now.Format("15:04")
	minuteLater := now.Add(1 * time.Minute).Format("15:04")
	today := now.Format("02-01-2006")
	filter := bson.M{
		"date": today,
		"checkTime": bson.M{
			"$gte": currentTime,
			"$lt":  minuteLater,
		},
	}

	cursor, err := s.collection.Find(ctx, filter, options.Find().SetSort(bson.D{
		{Key: "date", Value: 1},
		{Key: "checkTime", Value: 1},
	}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var schedule ScheduleLog
		if err := cursor.Decode(&schedule); err != nil {
			return nil, err
		}
		groupIds = append(groupIds, schedule.GroupId)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return groupIds, nil
}

func (s *scheduleLogRepositoryDB) InsertScheduleLog(scheduleLog *ScheduleLog) error {
	ctx := context.Background()
	_, err := s.collection.InsertOne(ctx, scheduleLog)
	return err
}

func (r *scheduleLogRepositoryDB) BatchInsertScheduleLogs(schedules []ScheduleLog) error {
	ctx := context.Background()
	var docs []interface{}
	for _, schedule := range schedules {
		docs = append(docs, schedule)
	}

	opts := options.InsertMany().SetOrdered(false)

	_, err := r.collection.InsertMany(ctx, docs, opts)
	if err != nil {
		return fmt.Errorf("failed to batch insert schedules: %v", err)
	}
	return nil
}

func (s *scheduleLogRepositoryDB) DeleteScheduleLog(groupId int) error {
	ctx := context.Background()
	filter := bson.M{"groupId": groupId}
	_, err := s.collection.DeleteMany(ctx, filter)
	return err
}

func (s *scheduleLogRepositoryDB) DeleteScheduleLogByRecurrenceId(recurrenceId int) error {
	ctx := context.Background()
	filter := bson.M{"recurrenceId": recurrenceId}
	_, err := s.collection.DeleteMany(ctx, filter)
	return err
}
