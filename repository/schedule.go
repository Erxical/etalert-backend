package repository

type Schedule struct {
	Id              string  `bson:"_id,omitempty"`
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

	Recurrence      string  `bson:"recurrence"`
	RecurrenceId    int     `bson:"recurrenceId"`
}

type Counter struct {
	ID  string `bson:"_id,omitempty"`
	Seq int    `bson:"seq"`
}

type ScheduleRepository interface {
	GetTravelTime(oriLat string, oriLong string, destLat string, destLong string, depTime string) (string, error)
	GetNextGroupId() (int, error)
	GetNextRecurrenceId() (int, error)
	CalculateNextRecurrenceDate(currentDate, recurrence string, count int) ([]string, error)
	BatchInsertSchedules(schedules []Schedule) error
	InsertSchedule(schedule *Schedule) error
	GetAllSchedules(gId string, date string) ([]*Schedule, error)
	GetScheduleById(id string) (*Schedule, error)
	GetSchedulesByGroupId(groupId int) ([]*Schedule, error)
	UpdateSchedule(id string, schedule *Schedule) error
	UpdateScheduleTime(id string, startTime string, endTime string) error
	DeleteSchedule(groupId int) error
	DeleteScheduleByRecurrenceId(recurrenceId int) error
}
