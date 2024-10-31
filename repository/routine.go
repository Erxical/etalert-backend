package repository

type Routine struct {
	Id       string   `bson:"_id,omitempty"`
	GoogleId string   `bson:"googleId"`
	Name     string   `bson:"name"`
	Duration int      `bson:"duration"`
	Order    int      `bson:"order"`
}

type RoutineRepository interface {
	InsertRoutine(routine *Routine) error
	GetAllRoutines(string) ([]*Routine, error)
	UpdateRoutine(string, *Routine) error
	DeleteRoutine(string) error
}
