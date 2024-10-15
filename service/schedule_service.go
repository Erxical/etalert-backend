package service

import (
	"etalert-backend/repository"
	"etalert-backend/websocket"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
)

type scheduleService struct {
	scheduleRepo    repository.ScheduleRepository
	scheduleLogRepo repository.ScheduleLogRepository
	routineRepo     repository.RoutineRepository
	bedtimeRepo     repository.BedtimeRepository
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, scheduleLogRepo repository.ScheduleLogRepository, routineRepo repository.RoutineRepository, bedTimeRepo repository.BedtimeRepository) ScheduleService {
	return &scheduleService{scheduleRepo: scheduleRepo, scheduleLogRepo: scheduleLogRepo, routineRepo: routineRepo, bedtimeRepo: bedTimeRepo}
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
		fmt.Println("Checking upcoming schedules...")
		s.autoUpdateSchedules()
	})
	c.Start()
}

func (s *scheduleService) autoUpdateSchedules() {
	groupIds, err := s.scheduleLogRepo.GetUpcomingSchedules()
	if err != nil {
		log.Printf("Failed to get upcoming schedules: %v", err)
	}
	for _, groupId := range groupIds {
		schedules, err := s.scheduleRepo.GetSchedulesByGroupId(groupId)
		newStartTime := schedules[0].StartTime
		newEndTime := newStartTime
		schedules = schedules[1:]
		if err != nil {
			log.Printf("Failed to get schedules by group ID: %v", err)
		}
		for _, schedule := range schedules {
			if schedule.IsUpdated {
				break
			}
			if schedule.IsTraveling {
				travelTimeText, err := s.scheduleRepo.GetTravelTime(
					fmt.Sprintf("%f", schedule.OriLatitude),
					fmt.Sprintf("%f", schedule.OriLongitude),
					fmt.Sprintf("%f", schedule.DestLatitude),
					fmt.Sprintf("%f", schedule.DestLongitude),
					"now",
				)
				if err != nil {
					log.Printf("Failed to get travel time: %v", err)
					return
				}

				startTime, err := time.Parse("15:04", newStartTime)
				if err != nil {
					log.Printf("failed to parse start time: %v", err)
				}

				travelDuration, err := parseDuration(travelTimeText)
				if err != nil {
					log.Printf("failed to parse travel duration: %v", err)
				}

				newEndTime = newStartTime
				newStartTime = startTime.Add(-travelDuration).Format("15:04")

				err = s.scheduleRepo.UpdateScheduleTime(schedule.Id, newStartTime, newEndTime)
				if err != nil {
					log.Printf("Failed to update schedule time: %v", err)
				}
				log.Printf("Updated schedule time for %s from user %s", schedule.Name, schedule.GoogleId)
			} else {
				startTime, err := time.Parse("15:04", schedule.StartTime)
				if err != nil {
					log.Printf("failed to parse start time: %v", err)
				}
				endTime, err := time.Parse("15:04", schedule.EndTime)
				if err != nil {
					log.Printf("failed to parse end time: %v", err)
				}
				duration := endTime.Sub(startTime)

				newEndTime = newStartTime
				newStartTime = startTime.Add(-duration).Format("15:04")

				err = s.scheduleRepo.UpdateScheduleTime(schedule.Id, newStartTime, newEndTime)
				if err != nil {
					log.Printf("Failed to update schedule time: %v", err)
				}
				log.Printf("Updated schedule time for %s from user %s", schedule.Name, schedule.GoogleId)
			}
		}
	}
}

func (s *scheduleService) InsertSchedule(schedule *ScheduleInput) (string, error) {
	groupId, err := s.scheduleRepo.GetNextGroupId()
	if err != nil {
		return "", fmt.Errorf("failed to get next group ID: %v", err)
	}
	schedule.GroupId = groupId

	err = s.scheduleRepo.InsertSchedule(&repository.Schedule{
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
		GroupId:         schedule.GroupId,
		Priority:        schedule.Priority,
		IsHaveLocation:  schedule.IsHaveLocation,
		IsFirstSchedule: schedule.IsFirstSchedule,
		IsTraveling:     schedule.IsTraveling,
		IsUpdated:       false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to insert schedule: %v", err)
	}

	if schedule.IsHaveLocation {
		err := s.handleTravelSchedule(schedule)
		if err != nil {
			log.Printf("Failed to handle travel schedule: %v", err)
		}
	}

	if schedule.IsFirstSchedule {
		str, err := s.insertRoutineSchedules(schedule)
		if err != nil {
			log.Printf("Failed to insert routines: %v", err)
		} else {
			return str, nil
		}
	}

	return "", nil
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
	leaveTime := startTime.Add(-travelDuration)

	leaveSchedule := &repository.Schedule{
		GoogleId:        schedule.GoogleId,
		Name:            "Leave From " + schedule.OriName,
		Date:            schedule.Date,
		StartTime:       leaveTime.Format("15:04"),
		EndTime:         schedule.StartTime,
		IsHaveEndTime:   false,
		OriName:         schedule.OriName,
		OriLatitude:     schedule.OriLatitude,
		OriLongitude:    schedule.OriLongitude,
		DestName:        schedule.DestName,
		DestLatitude:    schedule.DestLatitude,
		DestLongitude:   schedule.DestLongitude,
		GroupId:         schedule.GroupId,
		IsHaveLocation:  false,
		IsFirstSchedule: false,
		IsTraveling:     true,
		IsUpdated:       false,
	}

	err = s.scheduleRepo.InsertSchedule(leaveSchedule)
	if err != nil {
		return fmt.Errorf("failed to insert leave home schedule: %v", err)
	}

	scheduleLog := &repository.ScheduleLog{
		GroupId:       schedule.GroupId,
		OriLatitude:   schedule.OriLatitude,
		OriLongitude:  schedule.OriLongitude,
		DestLatitude:  schedule.DestLatitude,
		DestLongitude: schedule.DestLongitude,
		Date:          schedule.Date,
		CheckTime:     leaveTime.Add(-travelDuration).Format("15:04"),
	}

	err = s.scheduleLogRepo.InsertScheduleLog(scheduleLog)
	if err != nil {
		return fmt.Errorf("failed to insert schedule log: %v", err)
	}

	return nil
}

func (s *scheduleService) insertRoutineSchedules(schedule *ScheduleInput) (string, error) {
	if schedule == nil {
		return "", fmt.Errorf("schedule cannot be nil")
	}
	firstStartTime, err := s.scheduleRepo.GetFirstSchedule(schedule.GoogleId, schedule.Date)
	if err != nil {
		return "", fmt.Errorf("failed to get first schedule start time: %v", err)
	}

	routines, err := s.routineRepo.GetAllRoutines(schedule.GoogleId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user routines: %v", err)
	}

	currentStartTime, err := time.Parse("15:04", firstStartTime)
	if err != nil {
		return "", fmt.Errorf("failed to parse first schedule start time: %v", err)
	}

	for i := len(routines) - 1; i >= 0; i-- {
		routine := routines[i]
		routineDuration, err := parseDuration(fmt.Sprintf("%d min", routine.Duration))
		if err != nil {
			return "", fmt.Errorf("failed to parse routine duration: %v", err)
		}

		currentEndTime := currentStartTime
		currentStartTime = currentStartTime.Add(-routineDuration)

		newRoutineSchedule := &repository.Schedule{
			GoogleId:        schedule.GoogleId,
			Name:            routine.Name,
			Date:            schedule.Date,
			StartTime:       currentStartTime.Format("15:04"),
			EndTime:         currentEndTime.Format("15:04"),
			GroupId:         schedule.GroupId,
			IsHaveEndTime:   true,
			IsHaveLocation:  false,
			IsFirstSchedule: false,
			IsTraveling:     false,
			IsUpdated:       false,
		}

		err = s.scheduleRepo.InsertSchedule(newRoutineSchedule)
		if err != nil {
			log.Printf("Failed to insert routine schedule: %v", err) // Log and continue
			return "", fmt.Errorf("failed to insert routine schedule: %v", err)
		}
	}

	bedtimeStartTime := currentStartTime.Add(-5 * time.Minute)

	predefinedBedtime, err := s.bedtimeRepo.GetBedtimeInfo(schedule.GoogleId)
	if err != nil {
		return "", fmt.Errorf("failed to get predefined bedtime: %v", err)
	}

	// Parse the predefined bedtime for comparison
	predefinedBedtimeTime, err := time.Parse("15:04", predefinedBedtime.WakeTime)
	if err != nil {
		return "", fmt.Errorf("failed to parse predefined bedtime: %v", err)
	}

	bedtimeSchedule := &repository.Schedule{
		GoogleId:        schedule.GoogleId,
		Name:            "Wake up",
		Date:            schedule.Date,
		StartTime:       bedtimeStartTime.Format("15:04"),
		GroupId:         schedule.GroupId,
		IsHaveEndTime:   false,
		IsHaveLocation:  false,
		IsFirstSchedule: false,
		IsTraveling:     false,
		IsUpdated:       false,
	}

	err = s.scheduleRepo.InsertSchedule(bedtimeSchedule)
	if err != nil {
		log.Printf("Failed to insert bedtime schedule: %v", err) // Log and continue
		return "", fmt.Errorf("failed to insert bedtime schedule: %v", err)
	}

	// Compare the predefined bedtime with the auto-calculated bedtime
	if bedtimeStartTime.Before(predefinedBedtimeTime) {
		return "(auto-calculated bedtime is earlier than the predefined bedtime)", nil
	}

	return "", nil
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
			OriLatitude:     schedule.OriLatitude,
			OriLongitude:    schedule.OriLongitude,
			DestName:        schedule.DestName,
			DestLatitude:    schedule.DestLatitude,
			DestLongitude:   schedule.DestLongitude,
			GroupId:         schedule.GroupId,
			Priority:        schedule.Priority,
			IsHaveLocation:  schedule.IsHaveLocation,
			IsFirstSchedule: schedule.IsFirstSchedule,
			IsTraveling:     schedule.IsTraveling,
			IsUpdated:       schedule.IsUpdated,
		})
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
		GroupId:         schedule.GroupId,
		Priority:        schedule.Priority,
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

	updateMessage := []byte(fmt.Sprintf("Schedule updated for ID: %s", id))
	websocket.SendUpdate(updateMessage)

	return nil
}

func (s *scheduleService) DeleteSchedule(groupId string) error {
	id, err := strconv.Atoi(groupId)
	if err != nil {
		return fmt.Errorf("invalid groupId: %v", err)
	}
	err = s.scheduleRepo.DeleteSchedule(id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %v", err)
	}

	err = s.scheduleLogRepo.DeleteScheduleLog(id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule log: %v", err)
	}
	return nil
}
