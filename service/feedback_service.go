package service

import (
	"etalert-backend/repository"
)

type feedbackService struct {
	feedbackRepo repository.FeedbackRepository
}

func NewFeedbackService(feedbackRepo repository.FeedbackRepository) FeedbackService {
	return &feedbackService{feedbackRepo: feedbackRepo}
}

func (f feedbackService) InsertFeedback(feedback *FeedbackInput) error {

	err := f.feedbackRepo.InsertFeedback(&repository.Feedback{
		GoogleId: feedback.GoogleId,
		Feedback: feedback.Feedback,
	})
	if err != nil {
		return err
	}
	return nil
}
