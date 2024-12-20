package service

import "etalert-backend/repository"

type WeeklyReportResponse struct {
	Id        string                           `bson:"_id,omitempty"`
	Name      string                           `bson:"name"`
	StartDate string                           `bson:"startDate"`
	EndDate   string                           `bson:"endDate"`
	Tag       string                           `bson:"tag"`
	Details   []*repository.WeeklyReportDetail `bson:"details"`
}

type WeeklyReportService interface {
	StartCronJob()
	GetWeeklyReports(googleId string, date string) ([]WeeklyReportResponse, error)
}
