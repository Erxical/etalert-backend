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
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository) ScheduleService {
	return &scheduleService{scheduleRepo: scheduleRepo}
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

	if !schedule.IsHaveLocation {
		return nil
	}

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
	fmt.Println(travelDuration)

	leaveTime := startTime.Add(-travelDuration).Format("15:04")
	fmt.Println(leaveTime)
	leaveSchedule := &repository.Schedule{
		GoogleId:       schedule.GoogleId,
		Name:           "Leave Home At",
		Date:           schedule.Date, 
		StartTime:      leaveTime,
		EndTime:        "",
		IsHaveEndTime:  false,
		IsHaveLocation: false,
		IsFirstSchedule: false,
	}

	err = s.scheduleRepo.InsertSchedule(leaveSchedule)
	if err != nil {
		return fmt.Errorf("failed to insert leave home schedule: %v", err)
	}

	return nil
}
