package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type weeklyReportRepositoryDB struct {
	collection *mongo.Collection
}

func NewWeeklyReportRepositoryDB(client *mongo.Client, dbName string, collName string) WeeklyReportRepository {
	collection := client.Database(dbName).Collection(collName)
	return &weeklyReportRepositoryDB{collection: collection}
}

func (w *weeklyReportRepositoryDB) InsertWeeklyReport(weeklyReport *WeeklyReport) error {
	ctx := context.Background()
	_, err := w.collection.InsertOne(ctx, weeklyReport)
	return err
}

func (w *weeklyReportRepositoryDB) GetWeeklyReports(googleId string, date string) ([]*WeeklyReport, error) {
	ctx := context.Background()
	cursor, err := w.collection.Find(ctx, bson.M{"googleId": googleId, "startDate": date}, options.Find().SetSort(bson.D{{Key: "startDate", Value: 1}}))
	if err != nil {
		return nil, err
	}

	var weeklyReports []*WeeklyReport
	for cursor.Next(ctx) {
		var weeklyReport WeeklyReport
		err := cursor.Decode(&weeklyReport)
		if err != nil {
			return nil, err
		}
		weeklyReports = append(weeklyReports, &weeklyReport)
	}

	return weeklyReports, nil
}
