package service

import (
	"etalert-backend/repository"
	"sort"
)

type routineService struct {
	routineRepo repository.RoutineRepository
}

func NewRoutineService(routineRepo repository.RoutineRepository) RoutineService {
	return &routineService{routineRepo: routineRepo}
}

func (s routineService) InsertRoutine(routine *RoutineInput) error {

	err := s.routineRepo.InsertRoutine(&repository.Routine{
		GoogleId: routine.GoogleId,
		Name:     routine.Name,
		Duration: routine.Duration,
		Order:    routine.Order,
		Days:     routine.Days,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *routineService) GetAllRoutines(gId string) ([]*RoutineResponse, error) {
	dayOrder := map[string]int{
		"Sunday":    0,
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
	}

	routines, err := s.routineRepo.GetAllRoutines(gId)
	if err != nil {
		return nil, err
	}

	// Create a slice to hold the response
	var routineResponses []*RoutineResponse

	// Iterate over the routines to map each one to a RoutineResponse
	for _, routine := range routines {
		sort.Slice(routine.Days, func(i, j int) bool {
			return dayOrder[routine.Days[i]] < dayOrder[routine.Days[j]]
		})
		routineResponses = append(routineResponses, &RoutineResponse{
			Id:       routine.Id,
			Name:     routine.Name,
			Duration: routine.Duration,
			Order:    routine.Order,
			Days:     routine.Days,
		})
	}

	return routineResponses, nil
}

func (s *routineService) UpdateRoutine(id string, routine *RoutineUpdateInput) error {
	err := s.routineRepo.UpdateRoutine(id, &repository.Routine{
		Name:     routine.Name,
		Duration: routine.Duration,
		Order:    routine.Order,
		Days:     routine.Days,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *routineService) DeleteRoutine(id string) error {
	err := s.routineRepo.DeleteRoutine(id)
	if err != nil {
		return err
	}
	return nil
}
