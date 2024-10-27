package service

import (
	"etalert-backend/repository"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type weeklyReportService struct {
	weeklyReportRepo     repository.WeeklyReportRepository
	userRepo             repository.UserRepository
	routineRepo          repository.RoutineRepository
	weeklyReportListRepo repository.WeeklyReportListRepository
	routineLogRepo       repository.RoutineLogRepository
}

func NewWeeklyReportService(weeklyReportRepo repository.WeeklyReportRepository, userRepo repository.UserRepository, routineRepo repository.RoutineRepository, weeklyReportListRepo repository.WeeklyReportListRepository, routineLogRepo repository.RoutineLogRepository) WeeklyReportService {
	return &weeklyReportService{weeklyReportRepo: weeklyReportRepo, userRepo: userRepo, routineRepo: routineRepo, weeklyReportListRepo: weeklyReportListRepo, routineLogRepo: routineLogRepo}
}

func (w *weeklyReportService) StartCronJob() {
	c := cron.New()
	c.AddFunc("@every 1m", w.generateWeeklyReport)
	c.Start()
}

func (w *weeklyReportService) generateWeeklyReport() {
	now := time.Now().UTC().Add(7 * time.Hour)
	if now.Weekday() == time.Monday  && now.Hour() == 0 && now.Minute() == 0 {
	users, err := w.userRepo.GetAllUsersId()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, user := range users {
		routines, err := w.routineRepo.GetAllRoutines(user)
		if err != nil {
			fmt.Println(err)
			return
		}

		aWeekAgo := now.AddDate(0, 0, -7)
		routineReports, err := w.routineLogRepo.GetRoutineLogs(user, aWeekAgo.Format("02-01-2006"))
		if err != nil {
			fmt.Println(err)
			return
		}

		if len((routines)) != 0 && len(routineReports) != 0 {
			w.weeklyReportListRepo.InsertWeeklyReportList(&repository.WeeklyReportList{
				GoogleId:  user,
				StartDate: aWeekAgo.Format("02-01-2006"),
				EndDate:   now.AddDate(0, 0, -1).Format("02-01-2006"),
			})

			var weeklyReportDetails []*repository.WeeklyReportDetail
			for _, routine := range routines {
				for _, routineReport := range routineReports {
					if routine.Id == routineReport.RoutineId {
						weeklyReportDetail := &repository.WeeklyReportDetail{
							Date:          routineReport.Date,
							StartTime:     routineReport.StartTime,
							EndTime:       routineReport.EndTime,
							ActualEndTime: routineReport.ActualEndTime,
							Skewness:      routineReport.Skewness,
						}
						weeklyReportDetails = append(weeklyReportDetails, weeklyReportDetail)
					}
				}

				weeklyReport := &repository.WeeklyReport{
					GoogleId:  user,
					Name:      routine.Name,
					StartDate: aWeekAgo.Format("02-01-2006"),
					EndDate:   now.AddDate(0, 0, -1).Format("02-01-2006"),
					Days:      routine.Days,
					Details:   weeklyReportDetails,
				}
				w.weeklyReportRepo.InsertWeeklyReport(weeklyReport)
			}
		}
		fmt.Printf("Weekly report generated for %s \n", user)
		}
	}
}

func (w *weeklyReportService) GetWeeklyReports(googleId string, date string) ([]WeeklyReportResponse, error) {
	weeklyReports, err := w.weeklyReportRepo.GetWeeklyReports(googleId, date)
	if err != nil {
		return nil, err
	}

	var weeklyReportResponses []WeeklyReportResponse

	for _, weeklyReport := range weeklyReports {
		weeklyReportResponses = append(weeklyReportResponses, WeeklyReportResponse{
			Id:        weeklyReport.Id,
			Name:      weeklyReport.Name,
			StartDate: weeklyReport.StartDate,
			EndDate:   weeklyReport.EndDate,
			Days:      weeklyReport.Days,
			Details:   weeklyReport.Details,
		})
	}

	return weeklyReportResponses, nil
}
