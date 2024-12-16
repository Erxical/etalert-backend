package service

type WeeklyReportListResponse struct {
	Id        string `bson:"_id,omitempty"`
	StartDate string `bson:"startDate"`
	EndDate   string `bson:"endDate"`
}

type WeeklyReportListService interface {
	GetWeeklyReportLists(googleId string) ([]*WeeklyReportListResponse, error)
}
