package service

type TagInput struct {
	GoogleId string   `bson:"googleId"`
	Name     string   `bson:"name"`
	Routines []string `bson:"routines"`
}

type TagResponse struct {
	Id       string   `bson:"_id,omitempty"`
	Name     string   `bson:"name"`
	Routines []string `bson:"routines"`
}

type TagUpdateInput struct {
	Name     string   `bson:"name"`
	Routines []string `bson:"routines"`
}

type TagService interface {
	InsertTag(tag *TagInput) error
	GetAllTags(gId string) ([]*TagResponse, error)
	GetRoutinesByTagId(id string) ([]*RoutineResponse, error)
	UpdateTag(id string, tag *TagUpdateInput) error
	DeleteTag(id string) error
}
