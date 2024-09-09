package service

import (
	"etalert-backend/repository"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type scheduleService struct {
	scheduleRepo repository.ScheduleRepository
	routineRepo  repository.RoutineRepository
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, routineRepo repository.RoutineRepository) ScheduleService {
	return &scheduleService{scheduleRepo: scheduleRepo, routineRepo: routineRepo}
}

func parseDuration(durationText string) (time.Duration, error) {
	var totalMinutes int

	// Regular expressions to match various formats
	hoursPattern := regexp.MustCompile(`(\d+)\s*hour`)
	minutesPattern := regexp.MustCompile(`(\d+)\s*min`)

	// Find hours, if present
	hoursMatch := hoursPattern.FindStringSubmatch(durationText)
	if len(hoursMatch) > 1 {
		hours, err := strconv.Atoi(hoursMatch[1])
		if err != nil {
			return 0, fmt.Errorf("failed to parse hours: %v", err)
		}
		totalMinutes += hours * 60
	}

	// Find minutes, if present
	minutesMatch := minutesPattern.FindStringSubmatch(durationText)
	if len(minutesMatch) > 1 {
		minutes, err := strconv.Atoi(minutesMatch[1])
		if err != nil {
			return 0, fmt.Errorf("failed to parse minutes: %v", err)
		}
		totalMinutes += minutes
	}

	// Return total duration in minutes
	return time.Duration(totalMinutes) * time.Minute, nil
}

func (s *scheduleService) InsertSchedule(schedule *ScheduleInput) error {
	err := s.scheduleRepo.InsertSchedule(&repository.Schedule{
		GoogleId:        schedule.GoogleId,
		Name:            schedule.Name,
		Date:            schedule.Date,
		StartTime:       schedule.StartTime,
		EndTime:         schedule.EndTime,
		IsHaveEndTime:   schedule.IsHaveEndTime,
		Latitude:        schedule.DestLatitude,
		Longitude:       schedule.DestLongitude,
		IsHaveLocation:  schedule.IsHaveLocation,
		IsFirstSchedule: schedule.IsFirstSchedule,
	})
	if err != nil {
		return err
	}

	if schedule.IsHaveLocation {
		departureTime := schedule.DepartTime
		if departureTime == "" {
			departureTime = "now"
		}

		travelTimeText, err := s.scheduleRepo.GetTravelTime(
			fmt.Sprintf("%f", schedule.OriLatitude),
			fmt.Sprintf("%f", schedule.OriLongitude),
			fmt.Sprintf("%f", schedule.DestLatitude),
			fmt.Sprintf("%f", schedule.DestLongitude),
			departureTime,
		)

		if err != nil {
			return fmt.Errorf("failed to get travel time: %v", err)
		}

		startTime, err := time.Parse("15:04", schedule.StartTime)
		if err != nil {
			return fmt.Errorf("failed to parse start time: %v", err)
		}

		travelDuration, err := parseDuration(travelTimeText)
		if err != nil {
			return fmt.Errorf("failed to parse travel duration: %v", err)
		}
		travelDuration += 15 * time.Minute

		leaveTime := startTime.Add(-travelDuration).Format("15:04")
		leaveSchedule := &repository.Schedule{
			GoogleId:        schedule.GoogleId,
			Name:            "Leave Current Location",
			Date:            schedule.Date,
			StartTime:       leaveTime,
			EndTime:         "",
			IsHaveEndTime:   false,
			IsHaveLocation:  false,
			IsFirstSchedule: false,
		}

		err = s.scheduleRepo.InsertSchedule(leaveSchedule)
		if err != nil {
			return fmt.Errorf("failed to insert leave home schedule: %v", err)
		}
	}

	if schedule.IsFirstSchedule {
		firstStartTime, err := s.scheduleRepo.GetFirstSchedule(schedule.GoogleId, schedule.Date)
		if err != nil {
			return fmt.Errorf("failed to get first schedule start time: %v", err)
		}

		routines, err := s.routineRepo.GetAllRoutines(schedule.GoogleId)
		if err != nil {
			return fmt.Errorf("failed to fetch user routines: %v", err)
		}

		currentStartTime, err := time.Parse("15:04", firstStartTime)
		if err != nil {
			return fmt.Errorf("failed to parse first schedule start time: %v", err)
		}

		// Iterate over each routine in reverse order to adjust start times correctly
		for i := len(routines) - 1; i >= 0; i-- {
			routine := routines[i]
			routineDuration, err := parseDuration(fmt.Sprintf("%d min", routine.Duration))
			if err != nil {
				return fmt.Errorf("failed to parse routine duration: %v", err)
			}

			// Adjust current start time by subtracting the routine duration
			currentEndTime := currentStartTime
			currentStartTime = currentStartTime.Add(-routineDuration)
			newRoutineSchedule := &repository.Schedule{
				GoogleId:        schedule.GoogleId,
				Name:            routine.Name,
				Date:            schedule.Date,
				StartTime:       currentStartTime.Format("15:04"),
				EndTime:         currentEndTime.Format("15:04"), // Adjust if needed based on routine
				IsHaveEndTime:   false,
				IsHaveLocation:  false,
				IsFirstSchedule: false,
			}

			// Insert each adjusted routine as a schedule
			err = s.scheduleRepo.InsertSchedule(newRoutineSchedule)
			if err != nil {
				return fmt.Errorf("failed to insert routine schedule: %v", err)
			}
		}
	}

	return nil
}

func (s *scheduleService) GetAllSchedules(gId string, date string) ([]*ScheduleResponse, error) {
	schedules, err := s.scheduleRepo.GetAllSchedules(gId, date)
	if err != nil {
		return nil, err
	}

	var scheduleResponses []*ScheduleResponse

	for _, schedule := range schedules {
		scheduleResponses = append(scheduleResponses, &ScheduleResponse{
			Id:              schedule.Id,
			Name:            schedule.Name,
			StartTime:       schedule.StartTime,
			EndTime:         schedule.EndTime,
			IsHaveEndTime:   schedule.IsHaveEndTime,
			Latitude:        schedule.Latitude,
			Longitude:       schedule.Longitude,
			IsHaveLocation:  schedule.IsHaveLocation,
			IsFirstSchedule: schedule.IsFirstSchedule,
		})
	}

	for i, j := 0, len(scheduleResponses)-1; i < j; i, j = i+1, j-1 {
		scheduleResponses[i], scheduleResponses[j] = scheduleResponses[j], scheduleResponses[i]
	}

	return scheduleResponses, nil
}

func (s *scheduleService) UpdateSchedule(id string, schedule *ScheduleUpdateInput) error {
	// idFormat := "{$oid:\"" + id + "\"}"
	err := s.scheduleRepo.UpdateSchedule(id, &repository.Schedule{
		Name:          schedule.Name,
		Date:          schedule.Date,
		StartTime:     schedule.StartTime,
		EndTime:       schedule.EndTime,
		IsHaveEndTime: schedule.IsHaveEndTime,
	})
	if err != nil {
		return err
	}

	return nil
}
