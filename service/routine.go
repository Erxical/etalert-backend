package service

type RoutineInput struct {
	GoogleId string `bson:"googleId"`
	Name     string `bson:"name"`
	Duration int    `bson:"duration"`
	Order    int    `bson:"order"`
}

type RoutineResponse struct {
	Name     string `bson:"name"`
	Duration int    `bson:"duration"`
	Order    int    `bson:"order"`
}

type RoutineService interface {
	InsertRoutine(routine *RoutineInput) error
	GetRoutine(string) (*RoutineResponse, error)
	UpdateRoutine(string, *RoutineResponse) error
}
