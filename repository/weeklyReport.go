package repository

type WeeklyReport struct {
	Id        string                `bson:"_id,omitempty"`
	GoogleId  string                `bson:"googleId"`
	Name      string                `bson:"name"`
	StartDate string                `bson:"startDate"`
	EndDate   string                `bson:"endDate"`
	Details   []*WeeklyReportDetail `bson:"details"`
}

type WeeklyReportDetail struct {
	Date          string `bson:"date"`
	StartTime     string `bson:"startTime"`
	EndTime       string `bson:"endTime"`
	ActualEndTime string `bson:"actualEndTime"`
	Skewness      int    `bson:"skewness"`
}

type WeeklyReportRepository interface {
	InsertWeeklyReport(weeklyReport *WeeklyReport) error
	GetWeeklyReports(googleId string, date string) ([]*WeeklyReport, error)
}
