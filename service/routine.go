package service

type RoutineInput struct {
	GoogleId string   `bson:"googleId"`
	Name     string   `bson:"name"`
	Duration int      `bson:"duration"`
	Order    int      `bson:"order"`
}

type RoutineResponse struct {
	Id       string   `bson:"_id,omitempty"`
	Name     string   `bson:"name"`
	Duration int      `bson:"duration"`
	Order    int      `bson:"order"`
}

type RoutineUpdateInput struct {
	Name     string   `bson:"name"`
	Duration int      `bson:"duration"`
	Order    int      `bson:"order"`
}

type RoutineService interface {
	InsertRoutine(routine *RoutineInput) error
	GetAllRoutines(string) ([]*RoutineResponse, error)
	UpdateRoutine(string, *RoutineUpdateInput) error
	DeleteRoutine(string) error
}
