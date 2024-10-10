package service

type ScheduleInput struct {
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
	GroupId         int     `bson:"groupId"`
	Priority        int     `bson:"priority"`
	IsHaveLocation  bool    `bson:"isHaveLocation"`
	IsFirstSchedule bool    `bson:"isFirstSchedule"`
	IsTraveling     bool    `bson:"isTraveling"`
	IsUpdated       bool    `bson:"isUpdated"`
	DepartTime      string  `bson:"departTime"`
}

type ScheduleResponse struct {
	Id              string  `bson:"_id"`
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
	GroupId         int     `bson:"groupId"`
	Priority        int     `bson:"priority"`
	IsHaveLocation  bool    `bson:"isHaveLocation"`
	IsFirstSchedule bool    `bson:"isFirstSchedule"`
	IsTraveling     bool    `bson:"isTraveling"`
	IsUpdated       bool    `bson:"isUpdated"`
}

type ScheduleUpdateInput struct {
	Name          string `bson:"name"`
	Date          string `bson:"date"`
	StartTime     string `bson:"startTime"`
	EndTime       string `bson:"endTime"`
	IsHaveEndTime bool   `bson:"isHaveEndTime"`
}

type ScheduleService interface {
	StartCronJob()
	InsertSchedule(schedule *ScheduleInput) (string, error)
	GetAllSchedules(gId string, date string) ([]*ScheduleResponse, error)
	GetScheduleById(id string) (*ScheduleResponse, error)
	UpdateSchedule(id string, schedule *ScheduleUpdateInput) error
	DeleteSchedule(groupId string) error
}
