package service

import (
	"etalert-backend/repository"
)

type routineLogService struct {
	routineLogRepo repository.RoutineLogRepository
}

func NewRoutineLogService(routineLogRepo repository.RoutineLogRepository) RoutineLogService {
	return &routineLogService{routineLogRepo: routineLogRepo}
}

func (r *routineLogService) InsertRoutineLog(routineLog *RoutineLogInput) error {
	routineLogRepo := &repository.RoutineLog{
		RoutineId:     routineLog.RoutineId,
		GoogleId:      routineLog.GoogleId,
		Date:          routineLog.Date,
		StartTime:     routineLog.StartTime,
		EndTime:       routineLog.EndTime,
		ActualEndTime: routineLog.ActualEndTime,
		Skewness:      routineLog.Skewness,
	}
	return r.routineLogRepo.InsertRoutineLog(routineLogRepo)
}

func (r *routineLogService) GetRoutineLogs(googleId string, date string) ([]*RoutineLogResponse, error) {
	routinesLogs, err := r.routineLogRepo.GetRoutineLogs(googleId, date)
	if err != nil {
		return nil, err
	}

	var routineLogResponses []*RoutineLogResponse

	for _, routineLog := range routinesLogs {
		routineLogResponses = append(routineLogResponses, &RoutineLogResponse{
			Id:            routineLog.Id,
			RoutineId:     routineLog.RoutineId,
			Date:          routineLog.Date,
			StartTime:     routineLog.StartTime,
			EndTime:       routineLog.EndTime,
			ActualEndTime: routineLog.ActualEndTime,
			Skewness:      routineLog.Skewness,
		})
	}

	return routineLogResponses, nil
}

func (r *routineLogService) DeleteRoutineLog(id string) error {
	err := r.routineLogRepo.DeleteRoutineLog(id)
	if err != nil {
		return err
	}

	return nil
}