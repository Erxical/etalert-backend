package service

import (
	"encoding/json"
	"etalert-backend/repository"
	"etalert-backend/websocket"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
				} else {
					updateMessage := map[string]interface{}{
						"id":            schedule.Id,
						"name":          schedule.Name,
						"date":          schedule.Date,
						"startTime":     newStartTime,
						"endTime":       newEndTime,
						"isHaveEndTime": schedule.IsHaveEndTime,
					}
					message, _ := json.Marshal(updateMessage)
					websocket.SendUpdate(message)
				}
				log.Printf("Updated schedule time for %s from user %s", schedule.Name, schedule.GoogleId)
			} else {
				startTime, err := time.Parse("15:04", schedule.StartTime)
				if err != nil {
					log.Printf("failed to parse start time: %v", err)
				}
				endTime := startTime.Add(5 * time.Minute)
				if schedule.EndTime != "" {
					endTime, err = time.Parse("15:04", schedule.EndTime)
					if err != nil {
						log.Printf("failed to parse end time: %v", err)
					}
				}
				duration := endTime.Sub(startTime)

				newEndTime = newStartTime
				startTimeForNext, _ := time.Parse("15:04", newStartTime)
				newStartTime = startTimeForNext.Add(-duration).Format("15:04")

				err = s.scheduleRepo.UpdateScheduleTime(schedule.Id, newStartTime, newEndTime)
				if err != nil {
					log.Printf("Failed to update schedule time: %v", err)
				} else {
					updateMessage := map[string]interface{}{
						"id":            schedule.Id,
						"name":          schedule.Name,
						"date":          schedule.Date,
						"startTime":     newStartTime,
						"endTime":       newEndTime,
						"isHaveEndTime": schedule.IsHaveEndTime,
					}
					message, _ := json.Marshal(updateMessage)
					websocket.SendUpdate(message)
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

	recurrenceId, err := s.scheduleRepo.GetNextRecurrenceId()
	if err != nil {
		return "", fmt.Errorf("failed to get next recurrence ID: %v", err)
	}
	schedule.RecurrenceId = recurrenceId

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

		Recurrence:   schedule.Recurrence,
		RecurrenceId: schedule.RecurrenceId,
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
	arriveTime := startTime
	if leaveTime.Year() < arriveTime.Year() {
		date, err := time.Parse("02-01-2006", schedule.Date)
		if err == nil {
			date = date.AddDate(0, 0, -1)
			schedule.Date = date.Format("02-01-2006")
		}
	}
	schedule.StartTime = leaveTime.Format("15:04")

	leaveSchedule := &repository.Schedule{
		GoogleId:        schedule.GoogleId,
		Name:            "Leave From " + schedule.OriName,
		Date:            schedule.Date,
		StartTime:       schedule.StartTime,
		EndTime:         arriveTime.Format("15:04"),
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
		RecurrenceId:    schedule.RecurrenceId,
	}

	err = s.scheduleRepo.InsertSchedule(leaveSchedule)
	if err != nil {
		return fmt.Errorf("failed to insert leave home schedule: %v", err)
	}

	scheduleLog := &repository.ScheduleLog{
		GroupId:       schedule.GroupId,
		RecurrenceId:  schedule.RecurrenceId,
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

	routines, err := s.routineRepo.GetAllRoutines(schedule.GoogleId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user routines: %v", err)
	}

	scheduleDate, err := time.Parse("02-01-2006", schedule.Date)
	if err != nil {
		return "", fmt.Errorf("failed to parse schedule date: %v", err)
	}
	scheduleDay := scheduleDate.Weekday().String()

	currentStartTime, err := time.Parse("15:04", schedule.StartTime)
	if err != nil {
		return "", fmt.Errorf("failed to parse first schedule start time: %v", err)
	}

	for i := len(routines) - 1; i >= 0; i-- {
		routine := routines[i]
		if containsDay(routine.Days, scheduleDay) {
			routineDuration, err := parseDuration(fmt.Sprintf("%d min", routine.Duration))
			if err != nil {
				return "", fmt.Errorf("failed to parse routine duration: %v", err)
			}

			currentEndTime := currentStartTime
			currentStartTime = currentStartTime.Add(-routineDuration)

			if currentStartTime.Year() < currentEndTime.Year() {
				date, err := time.Parse("02-01-2006", schedule.Date)
				if err == nil {
					date = date.AddDate(0, 0, -1)
					schedule.Date = date.Format("02-01-2006")
				}
			}

			newRoutineSchedule := &repository.Schedule{
				GoogleId:        schedule.GoogleId,
				RoutineId:       routine.Id,
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
				RecurrenceId:    schedule.RecurrenceId,
			}

			err = s.scheduleRepo.InsertSchedule(newRoutineSchedule)
			if err != nil {
				log.Printf("Failed to insert routine schedule: %v", err) // Log and continue
				return "", fmt.Errorf("failed to insert routine schedule: %v", err)
			}
		}
	}

	bedtimeStartTime := currentStartTime.Add(-5 * time.Minute)

	predefinedBedtime, err := s.bedtimeRepo.GetBedtimeInfo(schedule.GoogleId)
	if err != nil {
		return "", fmt.Errorf("failed to get predefined bedtime: %v", err)
	}

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
		RecurrenceId:    schedule.RecurrenceId,
	}

	err = s.scheduleRepo.InsertSchedule(bedtimeSchedule)
	if err != nil {
		log.Printf("Failed to insert bedtime schedule: %v", err) // Log and continue
		return "", fmt.Errorf("failed to insert bedtime schedule: %v", err)
	}

	if bedtimeStartTime.Before(predefinedBedtimeTime) {
		return "(auto-calculated bedtime is earlier than the predefined bedtime)", nil
	}

	return "", nil
}

func containsDay(days []string, day string) bool {
	for _, d := range days {
		if strings.EqualFold(d, day) {
			return true
		}
	}
	return false
}

func (s *scheduleService) InsertRecurrenceSchedule(schedule *ScheduleInput) (string, error) {
	recurrenceId, err := s.scheduleRepo.GetNextRecurrenceId()
	if err != nil {
		return "", fmt.Errorf("failed to get next recurrence ID: %v", err)
	}
	schedule.RecurrenceId = recurrenceId

	recurrenceCount := map[string]int{
		"daily":   365,
		"weekly":  52,
		"monthly": 12,
		"yearly":  5,
	}[schedule.Recurrence]

	if recurrenceCount == 0 {
		return "", fmt.Errorf("invalid recurrence type: %v", schedule.Recurrence)
	}

	dates, err := s.scheduleRepo.CalculateNextRecurrenceDate(schedule.Date, schedule.Recurrence, recurrenceCount)
	if err != nil {
		return "", err
	}

	var travelDuration time.Duration
	if schedule.IsHaveLocation {
		travelDuration, err = s.calculateTravelDurationOnce(schedule)
		if err != nil {
			return "", fmt.Errorf("failed to calculate travel duration: %v", err)
		}
	}

	routines, err := s.routineRepo.GetAllRoutines(schedule.GoogleId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user routines: %v", err)
	}

	const batchSize = 100
	allSchedules := make([]repository.Schedule, 0)
	scheduleLogs := make([]repository.ScheduleLog, 0, len(dates))

	for _, date := range dates {
		groupId, err := s.scheduleRepo.GetNextGroupId()
		if err != nil {
			return "", fmt.Errorf("failed to get next group ID: %v", err)
		}
		schedule.GroupId = groupId

		mainSchedule := repository.Schedule{
			GoogleId:        schedule.GoogleId,
			Name:            schedule.Name,
			Date:            date,
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
			Recurrence:      schedule.Recurrence,
			RecurrenceId:    schedule.RecurrenceId,
		}

		currentTime, err := time.Parse("15:04", mainSchedule.StartTime)
		if err != nil {
			return "", fmt.Errorf("failed to parse start time: %v", err)
		}

		dateSchedules := []repository.Schedule{mainSchedule}

		if schedule.IsHaveLocation {
			currentEndTime := currentTime
			currentTime = currentTime.Add(-travelDuration)

			if currentTime.Year() < currentEndTime.Year() {
				newDate, err := time.Parse("02-01-2006", date)
				if err == nil {
					newDate = newDate.AddDate(0, 0, -1)
					date = newDate.Format("02-01-2006")
				}
			}
			travelSchedule := repository.Schedule{
				GoogleId:        schedule.GoogleId,
				Name:            "Leave From " + schedule.OriName,
				Date:            date,
				StartTime:       currentTime.Format("15:04"),
				EndTime:         mainSchedule.StartTime,
				IsHaveEndTime:   true,
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
				RecurrenceId:    schedule.RecurrenceId,
			}
			dateSchedules = append([]repository.Schedule{travelSchedule}, dateSchedules...)

			scheduleLog := repository.ScheduleLog{
				GroupId:       schedule.GroupId,
				RecurrenceId:  schedule.RecurrenceId,
				OriLatitude:   schedule.OriLatitude,
				OriLongitude:  schedule.OriLongitude,
				DestLatitude:  schedule.DestLatitude,
				DestLongitude: schedule.DestLongitude,
				Date:          date,
				CheckTime:     currentTime.Format("15:04"),
			}
			scheduleLogs = append(scheduleLogs, scheduleLog)
		}

		if schedule.IsFirstSchedule {
			scheduleDate, err := time.Parse("02-01-2006", date)
			if err != nil {
				return "", fmt.Errorf("failed to parse schedule date: %v", err)
			}
			scheduleDay := scheduleDate.Weekday().String()

			for i := len(routines) - 1; i >= 0; i-- {
				routine := routines[i]
				if containsDay(routine.Days, scheduleDay) {
					routineDuration, err := parseDuration(fmt.Sprintf("%d min", routine.Duration))
					if err != nil {
						return "", fmt.Errorf("failed to parse routine duration: %v", err)
					}

					endTime := currentTime
					currentTime = currentTime.Add(-routineDuration)

					if currentTime.Year() < endTime.Year() {
						newDate, err := time.Parse("02-01-2006", date)
						if err == nil {
							newDate = newDate.AddDate(0, 0, -1)
							date = newDate.Format("02-01-2006")
						}
					}
					routineSchedule := repository.Schedule{
						GoogleId:        schedule.GoogleId,
						RoutineId:       routine.Id,
						Name:            routine.Name,
						Date:            date,
						StartTime:       currentTime.Format("15:04"),
						EndTime:         endTime.Format("15:04"),
						GroupId:         schedule.GroupId,
						IsHaveEndTime:   true,
						IsHaveLocation:  false,
						IsFirstSchedule: false,
						IsTraveling:     false,
						IsUpdated:       false,
						RecurrenceId:    schedule.RecurrenceId,
					}
					dateSchedules = append([]repository.Schedule{routineSchedule}, dateSchedules...)
				}
			}
			bedtimeStartTime := currentTime.Add(-5 * time.Minute)
			bedtimeSchedule := repository.Schedule{
				GoogleId:        schedule.GoogleId,
				Name:            "Wake up",
				Date:            date,
				StartTime:       bedtimeStartTime.Format("15:04"),
				GroupId:         schedule.GroupId,
				IsHaveEndTime:   false,
				IsHaveLocation:  false,
				IsFirstSchedule: false,
				IsTraveling:     false,
				IsUpdated:       false,
				RecurrenceId:    schedule.RecurrenceId,
			}

			dateSchedules = append([]repository.Schedule{bedtimeSchedule}, dateSchedules...)
		}

		allSchedules = append(allSchedules, dateSchedules...)
	}

	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(2)

	// Process main schedules
	go func() {
		defer wg.Done()
		for i := 0; i < len(allSchedules); i += batchSize {
			end := i + batchSize
			if end > len(allSchedules) {
				end = len(allSchedules)
			}
			if err := s.scheduleRepo.BatchInsertSchedules(allSchedules[i:end]); err != nil {
				errCh <- fmt.Errorf("failed to insert batch schedules: %v", err)
				return
			}
		}
	}()

	// Process travel schedules
	go func() {
		defer wg.Done()
		for i := 0; i < len(scheduleLogs); i += batchSize {
			end := i + batchSize
			if end > len(scheduleLogs) {
				end = len(scheduleLogs)
			}
			if err := s.scheduleLogRepo.BatchInsertScheduleLogs(scheduleLogs[i:end]); err != nil {
				errCh <- fmt.Errorf("failed to insert batch schedule logs: %v", err)
				return
			}
		}
	}()

	// Process schedule logs
	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err := <-errCh; err != nil {
		return "", err
	}

	return "", nil
}

func (s *scheduleService) calculateTravelDurationOnce(schedule *ScheduleInput) (time.Duration, error) {
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
		return 0, fmt.Errorf("failed to get travel time: %v", err)
	}

	return parseDuration(travelTimeText)
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
			RoutineId:       schedule.RoutineId,
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

			Recurrence:   schedule.Recurrence,
			RecurrenceId: schedule.RecurrenceId,
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
		RoutineId:       schedule.RoutineId,
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

		Recurrence:   schedule.Recurrence,
		RecurrenceId: schedule.RecurrenceId,
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
		Id:              currentSchedule.Id,
		RoutineId:       currentSchedule.RoutineId,
		GoogleId:        currentSchedule.GoogleId,
		Name:            schedule.Name,
		Date:            schedule.Date,
		StartTime:       schedule.StartTime,
		EndTime:         schedule.EndTime,
		IsHaveEndTime:   schedule.IsHaveEndTime,
		OriName:         currentSchedule.OriName,
		OriLatitude:     currentSchedule.OriLatitude,
		OriLongitude:    currentSchedule.OriLongitude,
		DestName:        currentSchedule.DestName,
		DestLatitude:    currentSchedule.DestLatitude,
		DestLongitude:   currentSchedule.DestLongitude,
		GroupId:         currentSchedule.GroupId,
		Priority:        currentSchedule.Priority,
		IsHaveLocation:  currentSchedule.IsHaveLocation,
		IsFirstSchedule: currentSchedule.IsFirstSchedule,
		IsTraveling:     currentSchedule.IsTraveling,
		IsUpdated:       false,
		Recurrence:      currentSchedule.Recurrence,
		RecurrenceId:    currentSchedule.RecurrenceId,
	}

	// Check if the start time has changed
	startTimeChanged := currentSchedule.StartTime != updatedSchedule.StartTime
	if startTimeChanged {
		// Fetch all schedules for the same date, ordered by StartTime
		allSchedules, err := s.scheduleRepo.GetSchedulesByGroupId(currentSchedule.GroupId)
		if err != nil {
			return fmt.Errorf("failed to fetch schedules for the day: %v", err)
		}

		// Adjust times from the updated schedule backward
		// Starting from the end of the list and moving toward the current schedule
		currentStartTime, err := time.Parse("15:04", updatedSchedule.StartTime)
		if err != nil {
			return fmt.Errorf("failed to parse new start time: %v", err)
		}

		date := updatedSchedule.Date

		for i := 1; i < len(allSchedules); i++ {
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
				duration = 5 * time.Minute
			}

			// Calculate the new end time as the current start time
			endTime := currentStartTime
			sch.EndTime = endTime.Format("15:04")

			// Adjust the current start time by subtracting the routine duration
			currentStartTime = currentStartTime.Add(-duration)

			if currentStartTime.Year() < endTime.Year() {
				newDate, err := time.Parse("02-01-2006", sch.Date)
				if err == nil {
					newDate = newDate.AddDate(0, 0, -1)
					date = newDate.Format("02-01-2006")
				}
			}

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
				// sch.EndTime = ""
			}

			// Update the schedule times
			sch.StartTime = currentStartTime.Format("15:04")

			// Save the adjusted schedule back to the database
			err = s.scheduleRepo.UpdateSchedule(sch.Id, &repository.Schedule{
				Id:              sch.Id,
				RoutineId:       sch.RoutineId,
				GoogleId:        sch.GoogleId,
				Name:            sch.Name,
				Date:            date,
				StartTime:       sch.StartTime,
				EndTime:         sch.EndTime,
				IsHaveEndTime:   sch.IsHaveEndTime,
				OriName:         sch.OriName,
				OriLatitude:     sch.OriLatitude,
				OriLongitude:    sch.OriLongitude,
				DestName:        sch.DestName,
				DestLatitude:    sch.DestLatitude,
				DestLongitude:   sch.DestLongitude,
				GroupId:         sch.GroupId,
				Priority:        sch.Priority,
				IsHaveLocation:  sch.IsHaveLocation,
				IsFirstSchedule: sch.IsFirstSchedule,
				IsTraveling:     sch.IsTraveling,
				IsUpdated:       false,
				Recurrence:      sch.Recurrence,
				RecurrenceId:    sch.RecurrenceId,
			})

			if err != nil {
				fmt.Printf("Failed to adjust schedule times for %s: %v\n", sch.Name, err)
				return fmt.Errorf("failed to adjust schedule times: %v", err)
			} else {
				fmt.Printf("Successfully updated schedule: %s\n", sch.Name)
			}
			updateMessage := map[string]interface{}{
				"id":            sch.Id,
				"name":          sch.Name,
				"date":          sch.Date,
				"startTime":     sch.StartTime,
				"endTime":       sch.EndTime,
				"isHaveEndTime": sch.IsHaveEndTime,
			}
			message, _ := json.Marshal(updateMessage)
			websocket.SendUpdate(message)
		}
	}

	// Update the primary schedule entry
	err = s.scheduleRepo.UpdateSchedule(id, updatedSchedule)
	if err != nil {
		return fmt.Errorf("failed to update schedule: %v", err)
	}

	updateMessage := updatedSchedule
	message, _ := json.Marshal(updateMessage)
	websocket.SendUpdate(message)

	return nil
}

func (s *scheduleService) UpdateScheduleByRecurrenceId(recurrenceId string, inputSchedule *ScheduleUpdateInput, date string) error {
	id, err := strconv.Atoi(recurrenceId)
	if err != nil {
		return fmt.Errorf("invalid recurrenceId: %v", err)
	}
	schedules, err := s.scheduleRepo.GetSchedulesByRecurrenceId(id, date)
	if err != nil {
		return fmt.Errorf("failed to get schedules by recurrence ID: %v", err)
	}

	newRecurrenceId, err := s.scheduleRepo.GetNextRecurrenceId()
	if err != nil {
		return fmt.Errorf("failed to get next recurrence ID: %v", err)
	}

	for _, schedule := range schedules {
		currentSchedule, err := s.scheduleRepo.GetScheduleById(schedule.Id)
		if err != nil {
			return fmt.Errorf("failed to fetch current schedule: %v", err)
		}

		inputDate, err := time.Parse("02-01-2006", inputSchedule.Date)
		if err != nil {
			return fmt.Errorf("failed to parse input schedule date: %v", err)
		}

		currentDate, err := time.Parse("02-01-2006", currentSchedule.Date)
		if err != nil {
			return fmt.Errorf("failed to parse current schedule date: %v", err)
		}

		finalDate := time.Date(currentDate.Year(), inputDate.Month(), inputDate.Day(), 0, 0, 0, 0, time.UTC)
		combinedDate := finalDate.Format("02-01-2006")

		updatedSchedule := &repository.Schedule{
			Id:              currentSchedule.Id,
			RoutineId:       currentSchedule.RoutineId,
			GoogleId:        currentSchedule.GoogleId,
			Name:            inputSchedule.Name,
			Date:            combinedDate,
			StartTime:       inputSchedule.StartTime,
			EndTime:         inputSchedule.EndTime,
			IsHaveEndTime:   inputSchedule.IsHaveEndTime,
			OriName:         currentSchedule.OriName,
			OriLatitude:     currentSchedule.OriLatitude,
			OriLongitude:    currentSchedule.OriLongitude,
			DestName:        currentSchedule.DestName,
			DestLatitude:    currentSchedule.DestLatitude,
			DestLongitude:   currentSchedule.DestLongitude,
			GroupId:         currentSchedule.GroupId,
			Priority:        currentSchedule.Priority,
			IsHaveLocation:  currentSchedule.IsHaveLocation,
			IsFirstSchedule: currentSchedule.IsFirstSchedule,
			IsTraveling:     currentSchedule.IsTraveling,
			IsUpdated:       false,
			Recurrence:      currentSchedule.Recurrence,
			RecurrenceId:    newRecurrenceId,
		}

		startTimeChanged := currentSchedule.StartTime != updatedSchedule.StartTime
		if startTimeChanged {
			allSchedules, err := s.scheduleRepo.GetSchedulesByGroupId(currentSchedule.GroupId)
			if err != nil {
				return fmt.Errorf("failed to fetch schedules for the day: %v", err)
			}

			currentStartTime, err := time.Parse("15:04", updatedSchedule.StartTime)
			if err != nil {
				return fmt.Errorf("failed to parse new start time: %v", err)
			}

			for i := 1; i < len(allSchedules); i++ {
				sch := allSchedules[i]

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
					duration = 5 * time.Minute
				}

				endTime := currentStartTime
				sch.EndTime = endTime.Format("15:04")

				currentStartTime = currentStartTime.Add(-duration)

				if currentStartTime.Year() < endTime.Year() {
					newDate, err := time.Parse("02-01-2006", sch.Date)
					if err == nil {
						newDate = newDate.AddDate(0, 0, -1)
						date = newDate.Format("02-01-2006")
					}
				}

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

				sch.StartTime = currentStartTime.Format("15:04")

				err = s.scheduleRepo.UpdateSchedule(sch.Id, &repository.Schedule{
					Id:              sch.Id,
					RoutineId:       sch.RoutineId,
					GoogleId:        sch.GoogleId,
					Name:            sch.Name,
					Date:            date,
					StartTime:       sch.StartTime,
					EndTime:         sch.EndTime,
					IsHaveEndTime:   sch.IsHaveEndTime,
					OriName:         sch.OriName,
					OriLatitude:     sch.OriLatitude,
					OriLongitude:    sch.OriLongitude,
					DestName:        sch.DestName,
					DestLatitude:    sch.DestLatitude,
					DestLongitude:   sch.DestLongitude,
					GroupId:         sch.GroupId,
					Priority:        sch.Priority,
					IsHaveLocation:  sch.IsHaveLocation,
					IsFirstSchedule: sch.IsFirstSchedule,
					IsTraveling:     sch.IsTraveling,
					IsUpdated:       false,
					Recurrence:      sch.Recurrence,
					RecurrenceId:    newRecurrenceId,
				})

				if err != nil {
					fmt.Printf("Failed to adjust schedule times for %s: %v\n", sch.Name, err)
					return fmt.Errorf("failed to adjust schedule times: %v", err)
				} else {
					fmt.Printf("Successfully updated schedule: %s\n", sch.Name)
				}
				updateMessage := map[string]interface{}{
					"id":            sch.Id,
					"name":          sch.Name,
					"date":          sch.Date,
					"startTime":     sch.StartTime,
					"endTime":       sch.EndTime,
					"isHaveEndTime": sch.IsHaveEndTime,
				}
				message, _ := json.Marshal(updateMessage)
				websocket.SendUpdate(message)
			}
		}

		err = s.scheduleRepo.UpdateSchedule(schedule.Id, updatedSchedule)
		if err != nil {
			return fmt.Errorf("failed to update schedule: %v", err)
		}

		updateMessage := updatedSchedule
		message, _ := json.Marshal(updateMessage)
		websocket.SendUpdate(message)

	}

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

func (s *scheduleService) DeleteScheduleByRecurrenceId(recurrenceId string, date string) error {
	id, err := strconv.Atoi(recurrenceId)
	if err != nil {
		return fmt.Errorf("invalid recurrenceId: %v", err)
	}
	err = s.scheduleRepo.DeleteScheduleByRecurrenceId(id, date)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %v", err)
	}

	err = s.scheduleLogRepo.DeleteScheduleLogByRecurrenceId(id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule log: %v", err)
	}
	return nil
}
