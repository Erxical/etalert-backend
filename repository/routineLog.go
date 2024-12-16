package repository

type RoutineLog struct {
	Id            string `bson:"_id,omitempty"`
	RoutineId     string `bson:"routineId"`
	GoogleId      string `bson:"googleId"`
	Date          string `bson:"date"`
	StartTime     string `bson:"startTime"`
	EndTime       string `bson:"endTime"`
	ActualEndTime string `bson:"actualEndTime"`
	Skewness      int    `bson:"skewness"`
}

type RoutineLogRepository interface {
	InsertRoutineLog(RoutineLog *RoutineLog) error
	GetRoutineLogs(googleId string, date string) ([]*RoutineLog, error)
	DeleteRoutineLog(id string) error
}
