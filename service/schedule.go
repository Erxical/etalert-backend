package service

type ScheduleInput struct {
	GoogleId        string  `bson:"googleId"`
	Name            string  `bson:"name"`
	Date            string  `bson:"date"`
	StartTime       string  `bson:"startTime"`
	EndTime         string  `bson:"endTime"`
	IsHaveEndTime   bool    `bson:"isHaveEndTime"`
	OriLatitude     float64 `bson:"oriLatitude"`
	OriLongitude    float64 `bson:"oriLongitude"`
	DestLatitude    float64 `bson:"destLatitude"`
	DestLongitude   float64 `bson:"destLongitude"`
	IsHaveLocation  bool    `bson:"isHaveLocation"`
	IsFirstSchedule bool    `bson:"isFirstSchedule"`
	DepartTime      string  `bson:"departTime"`
}

type ScheduleResponse struct {
	Id              string  `bson:"_id"`
	Name            string  `bson:"name"`
	StartTime       string  `bson:"startTime"`
	EndTime         string  `bson:"endTime"`
	IsHaveEndTime   bool    `bson:"isHaveEndTime"`
	Latitude        float64 `bson:"latitude"`
	Longitude       float64 `bson:"longitude"`
	IsHaveLocation  bool    `bson:"isHaveLocation"`
	IsFirstSchedule bool    `bson:"isFirstSchedule"`
}

type ScheduleUpdateInput struct {
	Name          string `bson:"name"`
	Date          string `bson:"date"`
	StartTime     string `bson:"startTime"`
	EndTime       string `bson:"endTime"`
	IsHaveEndTime bool   `bson:"isHaveEndTime"`
}

type ScheduleService interface {
	InsertSchedule(schedule *ScheduleInput) error
	GetAllSchedules(gId string, date string) ([]*ScheduleResponse, error)
	GetScheduleById(id string) (*ScheduleResponse, error)
	UpdateSchedule(id string, schedule *ScheduleUpdateInput) error
}
