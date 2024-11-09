package repository

type Tag struct {
	Id       string   `bson:"_id,omitempty"`
	GoogleId string   `bson:"googleId"`
	Name     string   `bson:"name"`
	Routines []string `bson:"routines"`
}

type TagRepository interface {
	InsertTag(tag *Tag) error
	GetAllTags(gId string) ([]*Tag, error)
	GetRoutinesByTagId(id string) ([]string, error)
	GetTagByRoutineId(string) (*Tag, error)
	UpdateTag(id string, tag *Tag) error
	DeleteTag(id string) error
}
