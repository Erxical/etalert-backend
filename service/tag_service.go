package service

import "etalert-backend/repository"

type tagService struct {
	tagRepo repository.TagRepository
	routineRepo repository.RoutineRepository
}

func NewTagService(tagRepo repository.TagRepository, routineRepo repository.RoutineRepository) TagService {
	return &tagService{tagRepo: tagRepo, routineRepo: routineRepo}
}

func (s tagService) InsertTag(tag *TagInput) error {
	err := s.tagRepo.InsertTag(&repository.Tag{
		GoogleId: tag.GoogleId,
		Name:     tag.Name,
		Routines: tag.Routines,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *tagService) GetAllTags(gId string) ([]*TagResponse, error) {
	tags, err := s.tagRepo.GetAllTags(gId)
	if err != nil {
		return nil, err
	}

	var tagResponses []*TagResponse

	for _, tag := range tags {
		tagResponses = append(tagResponses, &TagResponse{
			Id:       tag.Id,
			Name:     tag.Name,
			Routines: tag.Routines,
		})
	}

	return tagResponses, nil
}

func (s *tagService) GetRoutinesByTagId(id string) ([]*RoutineResponse, error) {
	routineList, err := s.tagRepo.GetRoutinesByTagId(id)
	if err != nil {
		return nil, err
	}

	var routineResponses []*RoutineResponse

	for _, routineId := range routineList {
		routine, err := s.routineRepo.GetRoutineById(routineId)
		if err != nil {
			return nil, err
		}
		routineResponses = append(routineResponses, &RoutineResponse{
			Id:       routine.Id,
			Name:     routine.Name,
			Duration: routine.Duration,
			Order:    routine.Order,
		})
	}

	return routineResponses, nil
}

func (s *tagService) UpdateTag(id string, tag *TagUpdateInput) error {
	err := s.tagRepo.UpdateTag(id, &repository.Tag{
		Name:     tag.Name,
		Routines: tag.Routines,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *tagService) DeleteTag(id string) error {
	err := s.tagRepo.DeleteTag(id)
	if err != nil {
		return err
	}
	return nil
}