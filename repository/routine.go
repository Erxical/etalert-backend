package repository

type Routine struct {
	GoogleId string `bson:"googleId"`
	Name     string `bson:"name"`
	Duration int    `bson:"duration"`
	Order    int    `bson:"order"`
}

type RoutineRepository interface {
	GetHighestOrder(string) (int, error)
	InsertRoutine(routine *Routine) error
	GetAllRoutines(string) ([]*Routine, error)
	UpdateRoutine(string, *Routine) error
}
