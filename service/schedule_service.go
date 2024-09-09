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
		DestLatitude:    schedule.DestLatitude,
		DestLongitude:   schedule.DestLongitude,
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
			Latitude:        schedule.DestLatitude,
			Longitude:       schedule.DestLongitude,
			IsHaveLocation:  schedule.IsHaveLocation,
			IsFirstSchedule: schedule.IsFirstSchedule,
		})
	}

	for i, j := 0, len(scheduleResponses)-1; i < j; i, j = i+1, j-1 {
		scheduleResponses[i], scheduleResponses[j] = scheduleResponses[j], scheduleResponses[i]
	}

	return scheduleResponses, nil
}

func (s *scheduleService) GetScheduleById(id string) (*ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.GetScheduleById(id)
	if err != nil {
		return nil, err
	}

	return &ScheduleResponse{
		Id:              schedule.Id,
		Name:            schedule.Name,
		StartTime:       schedule.StartTime,
		EndTime:         schedule.EndTime,
		IsHaveEndTime:   schedule.IsHaveEndTime,
		Latitude:        schedule.DestLatitude,
		Longitude:       schedule.DestLongitude,
		IsHaveLocation:  schedule.IsHaveLocation,
		IsFirstSchedule: schedule.IsFirstSchedule,
	}, nil
}

func (s *scheduleService) UpdateSchedule(id string, schedule *ScheduleUpdateInput) error {
	currentSchedule, err := s.scheduleRepo.GetScheduleById(id)
	if err != nil {
		return fmt.Errorf("failed to fetch current schedule: %v", err)
	}

	updatedSchedule := &repository.Schedule{
		Name:          schedule.Name,
		Date:          schedule.Date,
		StartTime:     schedule.StartTime,
		EndTime:       schedule.EndTime,
		IsHaveEndTime: schedule.IsHaveEndTime,
	}

	startTimeChanged := currentSchedule.StartTime != updatedSchedule.StartTime
	if startTimeChanged && currentSchedule.IsHaveLocation {
		// Recalculate travel time using the distance matrix API
		travelTimeText, err := s.scheduleRepo.GetTravelTime(
			fmt.Sprintf("%f", currentSchedule.OriLatitude),
			fmt.Sprintf("%f", currentSchedule.OriLongitude),
			fmt.Sprintf("%f", currentSchedule.DestLatitude),
			fmt.Sprintf("%f", currentSchedule.DestLongitude),
			"now", // Use appropriate departure time or set dynamically if needed
		)
		if err != nil {
			return fmt.Errorf("failed to recalculate travel time: %v", err)
		}

		// Parse the recalculated travel time
		travelDuration, err := parseDuration(travelTimeText)
		if err != nil {
			return fmt.Errorf("failed to parse travel duration: %v", err)
		}

		// Fetch all schedules for the same date
		allSchedules, err := s.scheduleRepo.GetAllSchedules(currentSchedule.GoogleId, currentSchedule.Date)
		if err != nil {
			return fmt.Errorf("failed to fetch schedules for the day: %v", err)
		}

		// Get the time object of the updated start time
		newStartTime, err := time.Parse("15:04", updatedSchedule.StartTime)
		if err != nil {
			return fmt.Errorf("failed to parse new start time: %v", err)
		}

		// Adjust the times of schedules that are before the updated schedule and have not passed
		for _, sch := range allSchedules {
			// Parse the current schedule's start and end times
			scheduleStartTime, err := time.Parse("15:04", sch.StartTime)
			if err != nil {
				continue // Skip if time parsing fails
			}
			scheduleEndTime, err := time.Parse("15:04", sch.EndTime)
			if err != nil {
				continue // Skip if time parsing fails
			}

			// Calculate the duration of the current schedule
			scheduleDuration := scheduleEndTime.Sub(scheduleStartTime)

			// Skip schedules that are after the updated schedule or have already passed
			if scheduleStartTime.After(newStartTime) || scheduleStartTime.Before(time.Now()) {
				continue
			}

			// Adjust the start and end times based on the recalculated travel duration
			sch.StartTime = newStartTime.Add(-travelDuration).Format("15:04")
			sch.EndTime = newStartTime.Add(-travelDuration).Add(scheduleDuration).Format("15:04") // Maintain the original duration
			newStartTime = newStartTime.Add(-travelDuration)                                      // Update for the next iteration

			// Update the schedule in the database
			err = s.scheduleRepo.UpdateSchedule(sch.Id, &repository.Schedule{
				Name:          sch.Name,
				Date:          sch.Date,
				StartTime:     sch.StartTime,
				EndTime:       sch.EndTime,
				IsHaveEndTime: sch.IsHaveEndTime,
			})

			if err != nil {
				return fmt.Errorf("failed to adjust schedule times: %v", err)
			}
		}
	}

	err = s.scheduleRepo.UpdateSchedule(id, updatedSchedule)
	if err != nil {
		return fmt.Errorf("failed to update schedule: %v", err)
	}

	return nil
}
