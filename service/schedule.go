package service

type ScheduleInput struct {
	GoogleId        string  `bson:"googleId"`
	RoutineId       string  `bson:"routineId"`
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
	Transportation  string  `bson:"transportation"`
	Priority        int     `bson:"priority"`
	IsHaveLocation  bool    `bson:"isHaveLocation"`
	IsFirstSchedule bool    `bson:"isFirstSchedule"`
	IsTraveling     bool    `bson:"isTraveling"`
	IsUpdated       bool    `bson:"isUpdated"`
	DepartTime      string  `bson:"departTime"`
	TagId           string  `bson:"tagId"`

	Recurrence      string  `bson:"recurrence"`
	RecurrenceId    int     `bson:"recurrenceId"`
}

type ScheduleResponse struct {
	Id              string  `bson:"_id"`
	RoutineId       string  `bson:"routineId"`
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
	Transportation  string  `bson:"transportation"`
	Priority        int     `bson:"priority"`
	IsHaveLocation  bool    `bson:"isHaveLocation"`
	IsFirstSchedule bool    `bson:"isFirstSchedule"`
	IsTraveling     bool    `bson:"isTraveling"`
	IsUpdated       bool    `bson:"isUpdated"`
	TagId           string  `bson:"tagId"`
	
	Recurrence      string  `bson:"recurrence"`
	RecurrenceId    int     `bson:"recurrenceId"`
}

type ScheduleUpdateInput struct {
	Name          string `bson:"name"`
	Date          string `bson:"date"`
	StartTime     string `bson:"startTime"`
	EndTime       string `bson:"endTime"`
	IsHaveEndTime bool   `bson:"isHaveEndTime"`
}

type Traffic struct {
	Description string `bson:"description"`
	Cause       string `bson:"cause"`
	FromRoad    string `bson:"fromRoad"`
	ToRoad      string `bson:"toRoad"`
}

type Weather struct {
	Hazard            int    `bson:"hazard"`
	Weather           string `bson:"weather"`
	PrecipitationType string `bson:"precipitationType"`
}

type ScheduleService interface {
	StartCronJob()
	GetTraffic(oriLat string, oriLong string, destLat string, destLong string) ([]Traffic, error)
	GetWeather(oriLat string, oriLong string, destLat string, destLong string, travelTime string) ([]Weather, error)
	InsertSchedule(schedule *ScheduleInput) (string, error)
	InsertRecurrenceSchedule(schedule *ScheduleInput) (string, error)
	GetAllSchedules(gId string, date string) ([]*ScheduleResponse, error)
	GetScheduleById(id string) (*ScheduleResponse, error)
	GetSchedulesByGroupId(groupId string) ([]string, error)
	GetSchedulesIdByRecurrenceId(recurrenceId string, date string) ([]string, error)
	UpdateSchedule(id string, schedule *ScheduleUpdateInput) error
	UpdateScheduleByRecurrenceId(recurrenceId string, schedule *ScheduleUpdateInput, date string) error
	DeleteSchedule(groupId string) error
	DeleteScheduleByRecurrenceId(recurrenceId string, date string) error
}
