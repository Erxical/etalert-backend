package service

type RoutineInput struct {
	GoogleId string   `bson:"googleId"`
	Name     string   `bson:"name"`
	Duration int      `bson:"duration"`
	Order    int      `bson:"order"`
	Days     []string `bson:"days"`
}

type RoutineResponse struct {
	Id       string   `bson:"_id,omitempty"`
	Name     string   `bson:"name"`
	Duration int      `bson:"duration"`
	Order    int      `bson:"order"`
	Days     []string `bson:"days"`
}

type RoutineUpdateInput struct {
	Name     string   `bson:"name"`
	Duration int      `bson:"duration"`
	Order    int      `bson:"order"`
	Days     []string `bson:"days"`
}

type RoutineService interface {
	InsertRoutine(routine *RoutineInput) error
	GetAllRoutines(string) ([]*RoutineResponse, error)
	UpdateRoutine(string, *RoutineUpdateInput) error
}
