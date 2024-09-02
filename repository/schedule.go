package repository

type Schedule struct {
	GoogleId string `bson:"googleId"`
	Name	 string `bson:"name"`
	Date	 string `bson:"date"`
	StartTime string `bson:"startTime"`
	EndTime string `bson:"endTime"`
	IsHaveEndTime bool `bson:"isHaveEndTime"`
	Latitude float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
	IsHaveLocation bool `bson:"isHaveLocation"`
	IsFirstSchedule bool `bson:"isFirstSchedule"`
}

type ScheduleRepository interface {
	GetFirstSchedule(string, string) (string, error)
	GetTravelTime(string, string, string, string, string) (string, error)
	InsertSchedule(schedule *Schedule) error
	GetAllSchedules(string, string) ([]*Schedule, error)
}