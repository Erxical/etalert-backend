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
	Priority        int     `bson:"priority"`
	IsHaveLocation  bool    `bson:"isHaveLocation"`
	IsFirstSchedule bool    `bson:"isFirstSchedule"`
	IsTraveling     bool    `bson:"isTraveling"`
	IsUpdated       bool    `bson:"isUpdated"`
	DepartTime      string  `bson:"departTime"`

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
	Priority        int     `bson:"priority"`
	IsHaveLocation  bool    `bson:"isHaveLocation"`
	IsFirstSchedule bool    `bson:"isFirstSchedule"`
	IsTraveling     bool    `bson:"isTraveling"`
	IsUpdated       bool    `bson:"isUpdated"`

	Recurrence      string  `bson:"recurrence"`
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
	InsertRecurrenceSchedule(schedule *ScheduleInput) (string, error)
	GetAllSchedules(gId string, date string) ([]*ScheduleResponse, error)
	GetScheduleById(id string) (*ScheduleResponse, error)
	UpdateSchedule(id string, schedule *ScheduleUpdateInput) error
	DeleteSchedule(groupId string) error
	DeleteScheduleByRecurrenceId(recurrenceId string) error
}
