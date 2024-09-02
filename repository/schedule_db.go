package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	// "go.mongodb.org/mongo-driver/bson"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
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

func (s *scheduleRepositoryDB) InsertSchedule(schedule *Schedule) error {
	ctx := context.Background()
	_, err := s.collection.InsertOne(ctx, schedule)
	return err
}