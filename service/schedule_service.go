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
		OriLatitude:     schedule.OriLatitude,
		OriLongitude:    schedule.OriLongitude,
		DestLatitude:    schedule.DestLatitude,
		DestLongitude:   schedule.DestLongitude,
		IsHaveLocation:  schedule.IsHaveLocation,
		IsFirstSchedule: schedule.IsFirstSchedule,
		IsTraveling:     schedule.IsTraveling,
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
			IsTraveling:     true,
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
				IsHaveEndTime:   true,
				IsHaveLocation:  false,
				IsFirstSchedule: false,
				IsTraveling:     false,
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
			IsTraveling:     schedule.IsTraveling,
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
		IsTraveling:     schedule.IsTraveling,
	}, nil
}

func (s *scheduleService) UpdateSchedule(id string, schedule *ScheduleUpdateInput) error {
	// Fetch the current schedule by ID
	currentSchedule, err := s.scheduleRepo.GetScheduleById(id)
	if err != nil {
		return fmt.Errorf("failed to fetch current schedule: %v", err)
	}

	// Prepare the updated schedule structure
	updatedSchedule := &repository.Schedule{
		Name:          schedule.Name,
		Date:          schedule.Date,
		StartTime:     schedule.StartTime,
		EndTime:       schedule.EndTime,
		IsHaveEndTime: schedule.IsHaveEndTime,
	}

	// Check if the start time has changed
	startTimeChanged := currentSchedule.StartTime != updatedSchedule.StartTime
	if startTimeChanged {
		// Fetch all schedules for the same date, ordered by StartTime
		allSchedules, err := s.scheduleRepo.GetAllSchedules(currentSchedule.GoogleId, currentSchedule.Date)
		if err != nil {
			return fmt.Errorf("failed to fetch schedules for the day: %v", err)
		}

		// Adjust times from the updated schedule backward
		// Starting from the end of the list and moving toward the current schedule
		currentStartTime, err := time.Parse("15:04", updatedSchedule.StartTime)
		if err != nil {
			return fmt.Errorf("failed to parse new start time: %v", err)
		}

		for i := 0; i <= len(allSchedules)-1; i++ {
			sch := allSchedules[i]

			// Determine the duration of the current schedule
			var duration time.Duration
			if sch.IsHaveEndTime && sch.EndTime != "" {
				startTime, err := time.Parse("15:04", sch.StartTime)
				if err != nil {
					return fmt.Errorf("failed to parse start time: %v", err)
				}
				endTime, err := time.Parse("15:04", sch.EndTime)
				if err != nil {
					return fmt.Errorf("failed to parse end time: %v", err)
				}
				duration = endTime.Sub(startTime)
			} else {
				// Default to a minimal duration if no end time is set
				duration = 0
			}

			// Calculate the new end time as the current start time
			sch.EndTime = currentStartTime.Format("15:04")

			// Adjust the current start time by subtracting the routine duration
			currentStartTime = currentStartTime.Add(-duration)

			if (sch.IsTraveling) {
				travelTimeText, err := s.scheduleRepo.GetTravelTime(
					fmt.Sprintf("%f", allSchedules[i-1].OriLatitude),
					fmt.Sprintf("%f", allSchedules[i-1].OriLongitude),
					fmt.Sprintf("%f", allSchedules[i-1].DestLatitude),
					fmt.Sprintf("%f", allSchedules[i-1].DestLongitude),
					"now",
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

				leaveTime := startTime.Add(-travelDuration).Add(duration)
				currentStartTime = leaveTime
				sch.EndTime = ""
			}

			// Update the schedule times
			sch.StartTime = currentStartTime.Format("15:04")

			// Save the adjusted schedule back to the database
			err = s.scheduleRepo.UpdateSchedule(sch.Id, &repository.Schedule{
				Name:          sch.Name,
				Date:          sch.Date,
				StartTime:     sch.StartTime,
				EndTime:       sch.EndTime,
				IsHaveEndTime: sch.IsHaveEndTime,
			})
			if err != nil {
				fmt.Printf("Failed to adjust schedule times for %s: %v\n", sch.Name, err)
				return fmt.Errorf("failed to adjust schedule times: %v", err)
			} else {
				fmt.Printf("Successfully updated schedule: %s\n", sch.Name)
			}
		}
	}

	// Update the primary schedule entry
	err = s.scheduleRepo.UpdateSchedule(id, updatedSchedule)
	if err != nil {
		return fmt.Errorf("failed to update schedule: %v", err)
	}

	return nil
}
