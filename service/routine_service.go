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
    // Retrieve routines from the repository
    routines, err := s.routineRepo.GetAllRoutines(gId)
    if err != nil {
        return nil, err
    }

    // Create a slice to hold the response
    var routineResponses []*RoutineResponse

    // Iterate over the routines to map each one to a RoutineResponse
    for _, routine := range routines {
        routineResponses = append(routineResponses, &RoutineResponse{
            Name:     routine.Name,
            Duration: routine.Duration,
            Order:    routine.Order,
        })
    }

    return routineResponses, nil
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
