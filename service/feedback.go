package service

type FeedbackInput struct {
	GoogleId string `bson:"googleId"`
	Feedback string `bson:"feedback"`
}

type FeedbackService interface {
	InsertFeedback(feedback *FeedbackInput) error
}