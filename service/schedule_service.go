package service

import (
	"etalert-backend/repository"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
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

func (s *scheduleService) StartCronJob() {
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		s.checkUpcomingSchedules()
		fmt.Println("Checking upcoming schedules...")
	})
	c.Start()
}

func (s *scheduleService) checkUpcomingSchedules() {
	currentTime := time.Now()
	nextHour := currentTime.Add(1 * time.Hour)

	// Query schedules that are traveling, have startTime within the next hour, and aren't yet updated
	schedules, err := s.scheduleRepo.GetUpcomingTravelSchedules(currentTime, nextHour)
	if err != nil {
		log.Printf("Error querying upcoming travel schedules: %v", err)
		return
	}

	// For each upcoming travel schedule, recalculate travel time
	for _, schedule := range schedules {
		if !schedule.IsUpdated { // Only update if not already updated
			err := s.recalculateAndUpdateSchedule(schedule)
			if err != nil {
				log.Printf("Failed to update schedule %v: %v", schedule.GoogleId, err)
			}
		}
	}
}

func (s *scheduleService) recalculateAndUpdateSchedule(schedule *repository.Schedule) error {
	startTime, err := time.Parse("15:04", schedule.StartTime) // Assuming time is in "HH:mm" format
	if err != nil {
		return fmt.Errorf("failed to parse start time: %v", err)
	}
	// Step 1: Recalculate the travel time for the current schedule
	travelTimeText, err := s.scheduleRepo.GetTravelTime(
		fmt.Sprintf("%f", schedule.OriLatitude),
		fmt.Sprintf("%f", schedule.OriLongitude),
		fmt.Sprintf("%f", schedule.DestLatitude),
		fmt.Sprintf("%f", schedule.DestLongitude),
		"now",
	)
	if err != nil {
		return fmt.Errorf("failed to recalculate travel time: %v", err)
	}

	// Step 2: Parse the recalculated travel duration
	travelDuration, err := parseDuration(travelTimeText)
	if err != nil {
		return fmt.Errorf("failed to parse travel duration: %v", err)
	}

	// Step 3: Adjust the current schedule's leave time based on the new travel duration
	newLeaveTime := startTime.Add(-travelDuration - 15*time.Minute) // Example adjustment

	// Step 5: Convert newLeaveTime back to string format if needed
	newLeaveTimeStr := newLeaveTime.Format("15:04")

	// Update the current schedule with the new leave time
	err = s.scheduleRepo.UpdateScheduleStartTime(schedule.GoogleId, newLeaveTimeStr)
	if err != nil {
		return fmt.Errorf("failed to update schedule start time: %v", err)
	}

	log.Printf("Updated schedule for %v, new leave time: %v", schedule.GoogleId, newLeaveTimeStr)

	// Step 6: Shift the previous schedule's endTime earlier
	err = s.adjustPreviousSchedules(schedule.GoogleId, schedule.Date, newLeaveTime)
	if err != nil {
		return fmt.Errorf("failed to adjust previous schedules: %v", err)
	}

	return nil
}

func (s *scheduleService) adjustPreviousSchedules(googleId string, date string, newStartTime time.Time) error {
	// Find the previous schedule that ends just before the current schedule
	previousSchedule, err := s.scheduleRepo.GetPreviousSchedule(googleId, date, newStartTime)
	if err != nil {
		return fmt.Errorf("failed to fetch previous schedule: %v", err)
	}

	// If there's no previous schedule, there's nothing to adjust
	if previousSchedule == nil {
		return nil
	}

	// Adjust the end time of the previous schedule to end right before the new start time
	newEndTime := newStartTime.Add(-5 * time.Minute) // Adjusting to end 5 minutes before the current schedule's start time
	err = s.scheduleRepo.UpdateScheduleEndTime(previousSchedule.GoogleId, newEndTime)
	if err != nil {
		return fmt.Errorf("failed to update previous schedule end time: %v", err)
	}

	// Log the adjustment
	log.Printf("Updated previous schedule for %v, new end time: %v", previousSchedule.GoogleId, newEndTime)

	// Step 6: Recursively adjust previous schedules if needed
	return s.adjustPreviousSchedules(previousSchedule.GoogleId, previousSchedule.Date, newEndTime)
}

func (s *scheduleService) InsertSchedule(schedule *ScheduleInput) error {
	err := s.scheduleRepo.InsertSchedule(&repository.Schedule{
		GoogleId:        schedule.GoogleId,
		Name:            schedule.Name,
		Date:            schedule.Date,
		StartTime:       schedule.StartTime,
		EndTime:         schedule.EndTime,
		IsHaveEndTime:   schedule.IsHaveEndTime,
		OriName:         schedule.OriName,
		OriLatitude:     schedule.OriLatitude,
		OriLongitude:    schedule.OriLongitude,
		DestName:        schedule.DestName,
		DestLatitude:    schedule.DestLatitude,
		DestLongitude:   schedule.DestLongitude,
		IsHaveLocation:  schedule.IsHaveLocation,
		IsFirstSchedule: schedule.IsFirstSchedule,
		IsTraveling:     schedule.IsTraveling,
		IsUpdated:       false,
	})
	if err != nil {
		return fmt.Errorf("failed to insert schedule: %v", err)
	}

	if schedule.IsHaveLocation {
		err := s.handleTravelSchedule(schedule)
		if err != nil {
			log.Printf("Failed to handle travel schedule: %v", err)
		}
	}

	if schedule.IsFirstSchedule {
		err := s.insertRoutineSchedules(schedule)
		if err != nil {
			log.Printf("Failed to insert routines: %v", err)
		}
	}

	return nil
}

func (s *scheduleService) handleTravelSchedule(schedule *ScheduleInput) error {
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
		log.Printf("Failed to get travel time: %v", err)
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
		EndTime:         schedule.StartTime,
		IsHaveEndTime:   false,
		OriName:         schedule.OriName,
		OriLatitude:     schedule.OriLatitude,
		OriLongitude:    schedule.OriLongitude,
		DestName:        schedule.DestName,
		DestLatitude:    schedule.DestLatitude,
		DestLongitude:   schedule.DestLongitude,
		IsHaveLocation:  false,
		IsFirstSchedule: false,
		IsTraveling:     true,
		IsUpdated:       false,
	}

	err = s.scheduleRepo.InsertSchedule(leaveSchedule)
	if err != nil {
		log.Printf("Failed to insert leave home schedule: %v", err)
		return fmt.Errorf("failed to insert leave home schedule: %v", err)
	}
	return nil
}

func (s *scheduleService) insertRoutineSchedules(schedule *ScheduleInput) error {
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

	for i := len(routines) - 1; i >= 0; i-- {
		routine := routines[i]
		routineDuration, err := parseDuration(fmt.Sprintf("%d min", routine.Duration))
		if err != nil {
			return fmt.Errorf("failed to parse routine duration: %v", err)
		}

		currentEndTime := currentStartTime
		currentStartTime = currentStartTime.Add(-routineDuration)

		newRoutineSchedule := &repository.Schedule{
			GoogleId:        schedule.GoogleId,
			Name:            routine.Name,
			Date:            schedule.Date,
			StartTime:       currentStartTime.Format("15:04"),
			EndTime:         currentEndTime.Format("15:04"),
			IsHaveEndTime:   true,
			IsHaveLocation:  false,
			IsFirstSchedule: false,
			IsTraveling:     false,
			IsUpdated:       false,
		}

		err = s.scheduleRepo.InsertSchedule(newRoutineSchedule)
		if err != nil {
			log.Printf("Failed to insert routine schedule: %v", err) // Log and continue
			return fmt.Errorf("failed to insert routine schedule: %v", err)
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
			Date:            schedule.Date,
			StartTime:       schedule.StartTime,
			EndTime:         schedule.EndTime,
			IsHaveEndTime:   schedule.IsHaveEndTime,
			OriName:         schedule.OriName,
			DestName:        schedule.DestName,
			Latitude:        schedule.DestLatitude,
			Longitude:       schedule.DestLongitude,
			IsHaveLocation:  schedule.IsHaveLocation,
			IsFirstSchedule: schedule.IsFirstSchedule,
			IsTraveling:     schedule.IsTraveling,
		})
	}

	// for i, j := 0, len(scheduleResponses)-1; i < j; i, j = i+1, j-1 {
	// 	scheduleResponses[i], scheduleResponses[j] = scheduleResponses[j], scheduleResponses[i]
	// }

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
		Date:            schedule.Date,
		StartTime:       schedule.StartTime,
		EndTime:         schedule.EndTime,
		IsHaveEndTime:   schedule.IsHaveEndTime,
		OriName:         schedule.OriName,
		DestName:        schedule.DestName,
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

			if sch.IsTraveling {
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
