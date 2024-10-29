package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type scheduleRepositoryDB struct {
	collection *mongo.Collection
}

func NewScheduleRepositoryDB(client *mongo.Client, dbName string, collName string) ScheduleRepository {
	collection := client.Database(dbName).Collection(collName)
	return &scheduleRepositoryDB{collection: collection}
}

type DistanceMatrixResponse struct {
	Rows []struct {
		Elements []struct {
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

func (s *scheduleRepositoryDB) GetTravelTime(oriLat string, oriLong string, destLat string, destLong string, depTime string) (string, error) {
	godotenv.Load()
	apiKey := os.Getenv("G_MAP_API_KEY")
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s,%s&destinations=%s,%s&mode=driving&departure_time=%s&key=%s", oriLat, oriLong, destLat, destLong, depTime, apiKey)

	client := &http.Client{Timeout: 10 * time.Second}

	response, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make request to Google API: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var matrixResponse DistanceMatrixResponse
	if err := json.Unmarshal(body, &matrixResponse); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if matrixResponse.Status != "OK" || len(matrixResponse.Rows) == 0 || len(matrixResponse.Rows[0].Elements) == 0 {
		return "", fmt.Errorf("invalid response from Google API: %v", matrixResponse.Status)
	}

	element := matrixResponse.Rows[0].Elements[0]
	if element.Status != "OK" {
		return "", fmt.Errorf("invalid element status: %v", element.Status)
	}

	// Return the duration text
	return element.Duration.Text, nil
}

func (s *scheduleRepositoryDB) GetNextGroupId() (int, error) {
	var counter Counter
	ctx := context.Background()
	err := s.collection.FindOneAndUpdate(ctx, bson.M{"_id": "groupId"}, bson.M{"$inc": bson.M{"seq": 1}}, options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)).Decode(&counter)

	if err != nil {
		return 0, fmt.Errorf("failed to get next groupId: %v", err)
	}

	return counter.Seq, nil
}

func (s *scheduleRepositoryDB) GetNextRecurrenceId() (int, error) {
	var counter Counter
	ctx := context.Background()
	err := s.collection.FindOneAndUpdate(ctx, bson.M{"_id": "recurrenceId"}, bson.M{"$inc": bson.M{"seq": 1}}, options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)).Decode(&counter)

	if err != nil {
		return 0, fmt.Errorf("failed to get next recurrenceId: %v", err)
	}

	return counter.Seq, nil
}

func (s *scheduleRepositoryDB) CalculateNextRecurrenceDate(currentDate, recurrence string, count int) ([]string, error) {
	layout := "02-01-2006"
	date, err := time.Parse(layout, currentDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start date: %v", err)
	}

	var dates []string
	for i := 0; i < count; i++ {
		var nextDate time.Time
		switch recurrence {
		case "daily":
			nextDate = date.AddDate(0, 0, i)
		case "weekly":
			nextDate = date.AddDate(0, 0, i*7)
		case "monthly":
			nextDate = date.AddDate(0, i, 0)
		case "yearly":
			nextDate = date.AddDate(i, 0, 0)
		default:
			return nil, fmt.Errorf("invalid recurrence type: %v", recurrence)
		}
		dates = append(dates, nextDate.Format(layout))
	}
	return dates, nil
}

func (r *scheduleRepositoryDB) BatchInsertSchedules(schedules []Schedule) error {
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

func (s *scheduleRepositoryDB) InsertSchedule(schedule *Schedule) error {
	ctx := context.Background()
	_, err := s.collection.InsertOne(ctx, schedule)
	return err
}

func (s *scheduleRepositoryDB) GetAllSchedules(gId string, date string) ([]*Schedule, error) {
	ctx := context.Background()
	var schedules []*Schedule

	filter := bson.M{"googleId": gId}
	if date != "" {
		filter["date"] = date
	}

	cursor, err := s.collection.Find(ctx, filter, options.Find().SetSort(bson.D{
		{Key: "date", Value: 1},
		{Key: "startTime", Value: 1},
	}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var schedule Schedule
		if err := cursor.Decode(&schedule); err != nil {
			return nil, err
		}
		schedules = append(schedules, &schedule)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (s *scheduleRepositoryDB) GetScheduleById(id string) (*Schedule, error) {
	ctx := context.Background()
	var schedule Schedule

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ID: %v", err)
	}
	filter := bson.M{"_id": objectId}

	err = s.collection.FindOne(ctx, filter).Decode(&schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve schedule: %v", err)
	}

	return &schedule, nil
}

func (s *scheduleRepositoryDB) GetSchedulesByGroupId(groupId int) ([]*Schedule, error) {
	ctx := context.Background()
	var schedules []*Schedule

	filter := bson.M{"groupId": groupId}

	cursor, err := s.collection.Find(ctx, filter, options.Find().SetSort(bson.D{
		{Key: "date", Value: 1},
		{Key: "startTime", Value: -1},
	}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var schedule Schedule
		if err := cursor.Decode(&schedule); err != nil {
			return nil, err
		}
		schedules = append(schedules, &schedule)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (s *scheduleRepositoryDB) UpdateSchedule(id string, schedule *Schedule) error {
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ID: %v", err)
	}
	filter := bson.M{"_id": (objectId)}

	update := bson.M{"$set": bson.M{
		"googleId":        schedule.GoogleId,
		"routineId":       schedule.RoutineId,
		"name":            schedule.Name,
		"date":            schedule.Date,
		"startTime":       schedule.StartTime,
		"endTime":         schedule.EndTime,
		"isHaveEndTime":   schedule.IsHaveEndTime,
		"oriName":         schedule.OriName,
		"oriLatitude":     schedule.OriLatitude,
		"oriLongitude":    schedule.OriLongitude,
		"destName":        schedule.DestName,
		"destLatitude":    schedule.DestLatitude,
		"destLongitude":   schedule.DestLongitude,
		"groupId":         schedule.GroupId,
		"priority":        schedule.Priority,
		"isHaveLocation":  schedule.IsHaveLocation,
		"isFirstSchedule": schedule.IsFirstSchedule,
		"isTraveling":     schedule.IsTraveling,
		"isUpdated":       false,
		"recurrence":      schedule.Recurrence,
		"recurrenceId":    schedule.RecurrenceId,
	},
	}

	_, err = s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *scheduleRepositoryDB) UpdateScheduleTime(id string, startTime string, endTime string) error {
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ID: %v", err)
	}

	filter := bson.M{
		"_id": (objectId),
	}

	update := bson.M{"$set": bson.M{
		"startTime": startTime,
		"endTime":   endTime,
		"isUpdated": true,
	},
	}

	_, err = s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *scheduleRepositoryDB) DeleteSchedule(groupId int) error {
	ctx := context.Background()
	filter := bson.M{"groupId": groupId}

	_, err := s.collection.DeleteMany(ctx, filter)
	return err
}

func (s *scheduleRepositoryDB) DeleteScheduleByRecurrenceId(recurrenceId int, date string) error {
	ctx := context.Background()
	filter := bson.M{"recurrenceId": recurrenceId}
	if date != "" {
		filter["date"] = bson.M{"$gte": date}
	}

	_, err := s.collection.DeleteMany(ctx, filter)
	return err
}
