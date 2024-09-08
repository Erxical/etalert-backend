package service

type ScheduleInput struct {
	GoogleId        string  `bson:"googleId"`
	Name            string  `bson:"name"`
	Date            string  `bson:"date"`
	StartTime       string  `bson:"startTime"`
	EndTime         string  `bson:"endTime"`
	IsHaveEndTime   bool    `bson:"isHaveEndTime"`
	OriLatitude     float64 `bson:"latitude"`
	OriLongitude    float64 `bson:"longitude"`
	DestLatitude    float64 `bson:"latitude"`
	DestLongitude   float64 `bson:"longitude"`
	IsHaveLocation  bool    `bson:"isHaveLocation"`
	IsFirstSchedule bool    `bson:"isFirstSchedule"`
	DepartTime      string  `bson:"departTime"`
}

type ScheduleResponse struct {
	Name string `bson:"name"`
	StartTime string `bson:"startTime"`
	EndTime string `bson:"endTime"`
	IsHaveEndTime bool `bson:"isHaveEndTime"`
	Latitude float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
	IsHaveLocation bool `bson:"isHaveLocation"`
	IsFirstSchedule bool `bson:"isFirstSchedule"`
}

type ScheduleService interface {
	InsertSchedule(schedule *ScheduleInput) error
	GetAllSchedules(string, string) ([]*ScheduleResponse, error)
}
