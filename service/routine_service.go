package service

import (
	"etalert-backend/repository"
)

type routineService struct {
	routineRepo repository.RoutineRepository
}

func NewRoutineService(routineRepo repository.RoutineRepository) RoutineService {
	return &routineService{routineRepo: routineRepo}
}

func (s routineService) InsertRoutine(routine *RoutineInput) error {

	highestOrder, err := s.routineRepo.GetHighestOrder(routine.GoogleId)
	if err != nil {
		return err
	}
	// Increment the order number for the new routine
	newOrder := highestOrder + 1

	err = s.routineRepo.InsertRoutine(&repository.Routine{
		GoogleId: routine.GoogleId,
		Name:     routine.Name,
		Duration: routine.Duration,
		Order:    newOrder,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *routineService) GetRoutine(gId string) (*RoutineResponse, error) {
	routine, err := s.routineRepo.GetRoutine(gId)
	if err != nil {
		return nil, err
	}
	return &RoutineResponse{
		Name:     routine.Name,
		Duration: routine.Duration,
		Order:    routine.Order,
	}, nil

}

func (s *routineService) UpdateRoutine(gId string, routine *RoutineResponse) error {
	err := s.routineRepo.UpdateRoutine(gId, &repository.Routine{
		Name:     routine.Name,
		Duration: routine.Duration,
		Order:    routine.Order,
	})
	if err != nil {
		return err
	}
	return nil
}
