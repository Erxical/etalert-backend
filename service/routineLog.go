package service

type RoutineLogInput struct {
	RoutineId     string `bson:"routineId"`
	GoogleId      string `bson:"googleId"`
	Date          string `bson:"date"`
	StartTime     string `bson:"startTime"`
	EndTime       string `bson:"endTime"`
	ActualEndTime string `bson:"actualEndTime"`
	Skewness      int    `bson:"skewness"`
}

type RoutineLogResponse struct {
	Id            string `bson:"_id,omitempty"`
	RoutineId     string `bson:"routineId"`
	Date          string `bson:"date"`
	StartTime     string `bson:"startTime"`
	EndTime       string `bson:"endTime"`
	ActualEndTime string `bson:"actualEndTime"`
	Skewness      int    `bson:"skewness"`
}

type RoutineLogService interface {
	InsertRoutineLog(routineLog *RoutineLogInput) error
	GetRoutineLogs(googleId string, date string) ([]*RoutineLogResponse, error)
	DeleteRoutineLog(id string) error
}
