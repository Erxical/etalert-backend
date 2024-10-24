package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type weeklyReportListRepositoryDB struct {
	collection *mongo.Collection
}

func NewWeeklyReportListRepositoryDB(client *mongo.Client, dbName string, collName string) WeeklyReportListRepository {
	collection := client.Database(dbName).Collection(collName)
	return &weeklyReportListRepositoryDB{collection: collection}
}

func (w *weeklyReportListRepositoryDB) InsertWeeklyReportList(weeklyReportList *WeeklyReportList) error {
	ctx := context.Background()
	_, err := w.collection.InsertOne(ctx, weeklyReportList)
	return err
}

func (w *weeklyReportListRepositoryDB) GetWeeklyReportLists(googleId string) ([]*WeeklyReportList, error) {
	ctx := context.Background()
	cursor, err := w.collection.Find(ctx, bson.M{"googleId": googleId})
	if err != nil {
		return nil, err
	}

	var weeklyReportLists []*WeeklyReportList
	for cursor.Next(ctx) {
		var weeklyReportList WeeklyReportList
		err := cursor.Decode(&weeklyReportList)
		if err != nil {
			return nil, err
		}
		weeklyReportLists = append(weeklyReportLists, &weeklyReportList)
	}

	return weeklyReportLists, nil
}