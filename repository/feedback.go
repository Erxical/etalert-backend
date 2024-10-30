package repository

type Feedback struct {
	Id              string    `bson:"_id,omitempty"`
	GoogleId        string    `bson:"googleId"`
	Feedback        string    `bson:"feedback"`
}

type FeedbackRepository interface {
	InsertFeedback(feedback *Feedback) error
}
