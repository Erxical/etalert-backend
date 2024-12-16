package service

import (
	"etalert-backend/repository"
)

type weeklyReportListService struct {
	weeklyReportListRepo repository.WeeklyReportListRepository
}

func NewWeeklyReportListService(weeklyReportListRepo repository.WeeklyReportListRepository) WeeklyReportListService {
	return &weeklyReportListService{weeklyReportListRepo: weeklyReportListRepo}
}

func (w *weeklyReportListService) GetWeeklyReportLists(googleId string) ([]*WeeklyReportListResponse, error) {
	weeklyReportLists, err := w.weeklyReportListRepo.GetWeeklyReportLists(googleId)
	if err != nil {
		return nil, err
	}

	var weeklyReportListResponses []*WeeklyReportListResponse
	for _, weeklyReportList := range weeklyReportLists {
		weeklyReportListResponses = append(weeklyReportListResponses, &WeeklyReportListResponse{
			Id:        weeklyReportList.Id,
			StartDate: weeklyReportList.StartDate,
			EndDate:   weeklyReportList.EndDate,
		})
	}

	return weeklyReportListResponses, nil
}
