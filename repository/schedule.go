package repository

import "time"

type Schedule struct {
	Id              string  `bson:"_id,omitempty"`
	GoogleId        string  `bson:"googleId"`
	Name            string  `bson:"name"`
	Date            string  `bson:"date"`
	StartTime       string  `bson:"startTime"`
	EndTime         string  `bson:"endTime"`
	IsHaveEndTime   bool    `bson:"isHaveEndTime"`
	OriName         string  `bson:"oriName"`
	OriLatitude     float64 `bson:"oriLatitude"`
	OriLongitude    float64 `bson:"oriLongitude"`
	DestName        string  `bson:"destName"`
	DestLatitude    float64 `bson:"destLatitude"`
	DestLongitude   float64 `bson:"destLongitude"`
	IsHaveLocation  bool    `bson:"isHaveLocation"`
	IsFirstSchedule bool    `bson:"isFirstSchedule"`
	IsTraveling     bool    `bson:"isTraveling"`
	IsUpdated       bool    `bson:"isUpdated"`
}

type ScheduleRepository interface {
	GetUpcomingTravelSchedules(nextHour time.Time) ([]*Schedule, error)
	UpdateScheduleStartTime(googleId string, newStartTime string) error
	GetPreviousSchedule(googleId string, date string, newStartTime time.Time) (*Schedule, error)
	UpdateScheduleEndTime(googleId string, newEndTime time.Time) error
	GetFirstSchedule(gId string, date string) (string, error)
	GetTravelTime(oriLat string, oriLong string, destLat string, destLong string, depTime string) (string, error)
	InsertSchedule(schedule *Schedule) error
	GetAllSchedules(gId string, date string) ([]*Schedule, error)
	GetScheduleById(id string) (*Schedule, error)
	UpdateSchedule(id string, schedule *Schedule) error
}
