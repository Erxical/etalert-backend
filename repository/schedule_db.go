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

func (s *scheduleRepositoryDB) GetUpcomingTravelSchedules(currentTime, nextHour time.Time) ([]*Schedule, error) {
	ctx := context.Background()
	var schedules []*Schedule

	filter := bson.M{
		"isTraveling": true,
		"isUpdated":   false,
		"startTime": bson.M{
			"$gte": currentTime.Format("15:04"),
			"$lte": nextHour.Format("15:04"),
		},
	}
	cursor, err := s.collection.Find(ctx, filter)
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

	return schedules, nil
}

func (s *scheduleRepositoryDB) GetTravelSchedule(googleId string, date string) (*Schedule, error) {
	ctx := context.Background()
	var travelSchedule Schedule

	// Query for the travel schedule based on the googleId, date, and isTraveling = true
	filter := bson.M{
		"googleId":    googleId,
		"date":        date,
		"isTraveling": true, // Filter to find the travel schedule
	}

	err := s.collection.FindOne(ctx, filter).Decode(&travelSchedule)
	if err == mongo.ErrNoDocuments {
		return nil, nil // No travel schedule found
	}
	if err != nil {
		return nil, err // Return other errors if found
	}

	return &travelSchedule, nil
}

func (s *scheduleRepositoryDB) UpdateTravelScheduleTimes(googleId string, newStartTime string, newEndTime string) error {
	ctx := context.Background()

	filter := bson.M{
		"googleId":    googleId,
		"isTraveling": true, // Make sure we're updating the travel schedule
		"isUpdated":   false,
	}
	update := bson.M{
		"$set": bson.M{
			"startTime": newStartTime,
			"endTime":   newEndTime,
			"isUpdated": true,
		},
	}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *scheduleRepositoryDB) UpdateScheduleStartTime(googleId string, newStartTime string) error {
	ctx := context.Background()

	filter := bson.M{
		"googleId":  googleId,
		"isUpdated": false,
	}
	update := bson.M{
		"$set": bson.M{
			"startTime": newStartTime,
			"isUpdated": true,
		},
	}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *scheduleRepositoryDB) GetPreviousSchedule(googleId string, date string, newStartTime time.Time) (*Schedule, error) {
	ctx := context.Background()
	var schedule Schedule

	// Find the previous schedule that ends before newStartTime
	filter := bson.M{
		"googleId": googleId,
		"date":     date,
		"endTime": bson.M{
			"$lt": newStartTime.Format("15:04"), // Previous schedule should end before the new start time
		},
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "endTime", Value: -1}}) // Sort to get the latest one

	err := s.collection.FindOne(ctx, filter, opts).Decode(&schedule)
	if err == mongo.ErrNoDocuments {
		return nil, nil // No previous schedule found
	}
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (s *scheduleRepositoryDB) UpdateScheduleEndTime(googleId string, scheduleId string, newEndTime time.Time) error {
	ctx := context.Background()

	filter := bson.M{
		"googleId": googleId,
		"_id":      scheduleId, // Ensure the correct schedule is updated
	}
	update := bson.M{
		"$set": bson.M{
			"endTime":   newEndTime.Format("15:04"),
			"isUpdated": true,
		},
	}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
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

func (s *scheduleRepositoryDB) GetFirstSchedule(googleId string, date string) (string, error) {
	ctx := context.Background()
	var schedule Schedule

	// Define filter to match the date and Google ID
	filter := bson.M{"googleId": googleId, "date": date}

	// Find the first schedule of the day sorted by StartTime
	opts := options.FindOne().SetSort(bson.D{{Key: "startTime", Value: 1}}) // Sort by startTime in ascending order

	err := s.collection.FindOne(ctx, filter, opts).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("no schedules found for the given date")
		}
		return "", fmt.Errorf("failed to retrieve first schedule: %v", err)
	}

	return schedule.StartTime, nil
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

func (s *scheduleRepositoryDB) UpdateSchedule(id string, schedule *Schedule) error {
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ID: %v", err)
	}
	filter := bson.M{"_id": (objectId)}

	update := bson.M{"$set": bson.M{
		"name":          schedule.Name,
		"date":          schedule.Date,
		"startTime":     schedule.StartTime,
		"endTime":       schedule.EndTime,
		"isHaveEndTime": schedule.IsHaveEndTime,
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