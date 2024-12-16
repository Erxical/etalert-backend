package service

import (
	"etalert-backend/repository"
)

type routineService struct {
	routineRepo repository.RoutineRepository
	tagRepo     repository.TagRepository
}

func NewRoutineService(routineRepo repository.RoutineRepository, tagRepo repository.TagRepository) RoutineService {
	return &routineService{routineRepo: routineRepo, tagRepo: tagRepo}
}

func (s routineService) InsertRoutine(routine *RoutineInput) error {
	err := s.routineRepo.InsertRoutine(&repository.Routine{
		GoogleId: routine.GoogleId,
		Name:     routine.Name,
		Duration: routine.Duration,
		Order:    routine.Order,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *routineService) GetAllRoutines(gId string) ([]*RoutineResponse, error) {
	routines, err := s.routineRepo.GetAllRoutines(gId)
	if err != nil {
		return nil, err
	}

	// Create a slice to hold the response
	var routineResponses []*RoutineResponse

	// Iterate over the routines to map each one to a RoutineResponse
	for _, routine := range routines {
		routineResponses = append(routineResponses, &RoutineResponse{
			Id:       routine.Id,
			Name:     routine.Name,
			Duration: routine.Duration,
			Order:    routine.Order,
		})
	}

	return routineResponses, nil
}

func (s *routineService) UpdateRoutine(id string, routine *RoutineUpdateInput) error {
	currentRoutine, err := s.routineRepo.GetRoutineById(id)
	if err != nil {
		return err
	}
	err = s.routineRepo.UpdateRoutine(id, &repository.Routine{
		Id:       currentRoutine.Id,
		GoogleId: currentRoutine.GoogleId,
		Name:     routine.Name,
		Duration: routine.Duration,
		Order:    routine.Order,
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

	tag, err := s.tagRepo.GetTagByRoutineId(id)
	if err != nil {
		return err
	}

	var updatedRoutines []string
	for _, routineID := range tag.Routines {
		if routineID != id {
			updatedRoutines = append(updatedRoutines, routineID)
		}
	}

	if len(updatedRoutines) == 0 {
		updatedRoutines = []string{}
	}

	err = s.tagRepo.UpdateTag(tag.Id, &repository.Tag{
		Name:     tag.Name,
		Routines: updatedRoutines,
	})

	return nil
}
