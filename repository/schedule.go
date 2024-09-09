package repository

type Schedule struct {
	Id              string  `bson:"_id,omitempty"`
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
	IsTraveling     bool    `bson:"isTraveling"`
}

type ScheduleRepository interface {
	GetFirstSchedule(gId string, date string) (string, error)
	GetTravelTime(oriLat string, oriLong string, destLat string, destLong string, depTime string) (string, error)
	InsertSchedule(schedule *Schedule) error
	GetAllSchedules(gId string, date string) ([]*Schedule, error)
	GetScheduleById(id string) (*Schedule, error)
	UpdateSchedule(id string, schedule *Schedule) error
}
