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

func (s *scheduleRepositoryDB) InsertSchedule(schedule *Schedule) error {
	ctx := context.Background()
	_, err := s.collection.InsertOne(ctx, schedule)
	return err
}

func (s *scheduleRepositoryDB) GetAllSchedules(gId string, date string) ([]*Schedule, error) {
	ctx := context.Background()
	var schedules []*Schedule
	filter := bson.M{"googleId": gId, "date": date}

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
		// Append the decoded routine to the slice
		schedules = append(schedules, &schedule)
	}

	// Check if any errors occurred during iteration
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

	// Define update to replace the existing schedule
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

// {"_id":{"$oid":"66d592da997f1d43a5d0f2e9"}}
