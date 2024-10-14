package repository

type ScheduleLog struct {
	Id            string  `bson:"_id,omitempty"`
	GroupId       int     `bson:"groupId"`
	OriLatitude   float64 `bson:"oriLatitude"`
	OriLongitude  float64 `bson:"oriLongitude"`
	DestLatitude  float64 `bson:"destLatitude"`
	DestLongitude float64 `bson:"destLongitude"`
	CheckTime     string  `bson:"checkTime"`
}

type ScheduleLogRepository interface {
	InsertScheduleLog(scheduleLog *ScheduleLog) error
}
