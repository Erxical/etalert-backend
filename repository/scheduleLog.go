package repository

type ScheduleLog struct {
	Id            string  `bson:"_id,omitempty"`
	GroupId       int     `bson:"groupId"`
	RecurrenceId  int     `bson:"recurrenceId"`
	OriLatitude   float64 `bson:"oriLatitude"`
	OriLongitude  float64 `bson:"oriLongitude"`
	DestLatitude  float64 `bson:"destLatitude"`
	DestLongitude float64 `bson:"destLongitude"`
	Date          string  `bson:"date"`
	CheckTime     string  `bson:"checkTime"`
}

type ScheduleLogRepository interface {
	GetUpcomingSchedules() ([]int, error)
	InsertScheduleLog(scheduleLog *ScheduleLog) error
	BatchInsertScheduleLogs(schedules []ScheduleLog) error
	DeleteScheduleLog(groupId int) error
	DeleteScheduleLogByRecurrenceId(recurrenceId int) error
}
