package repository

type WeeklyReportList struct {
	Id        string                `bson:"_id,omitempty"`
	GoogleId  string                `bson:"googleId"`
	StartDate string                `bson:"startDate"`
	EndDate   string                `bson:"endDate"`
}

type WeeklyReportListRepository interface {
	InsertWeeklyReportList(weeklyReportList *WeeklyReportList) error
	GetWeeklyReportLists(googleId string) ([]*WeeklyReportList, error)
}