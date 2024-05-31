package service

import (
	"errors"
	"etalert-backend/repository"
)

type bedtimeService struct {
	bedtimeRepo repository.BedtimeRepository
}

func NewBedtimeService(bedtimeRepo repository.BedtimeRepository) BedtimeService {
	return &bedtimeService{bedtimeRepo: bedtimeRepo}
}

var ErrBedtimeAlreadyExists = errors.New("bedtime already exists")

func (s bedtimeService) InsertBedtime(bedtime *BedtimeInput) error {
	existingBedtime, err := s.bedtimeRepo.GetBedtimeInfo(bedtime.GoogleId)
	if err != nil {
		return err
	}
	if existingBedtime != nil {
		return ErrBedtimeAlreadyExists
	}

	err = s.bedtimeRepo.InsertBedtime(&repository.Bedtime{
		GoogleId:  bedtime.GoogleId,
		SleepTime: bedtime.SleepTime,
		WakeTime:  bedtime.WakeTime,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s bedtimeService) GetBedtimeInfo(gId string) (*BedtimeResponse, error) {
	bedtime, err := s.bedtimeRepo.GetBedtimeInfo(gId)
	if err != nil {
		return nil, err
	}

	bedtimeResponse := BedtimeResponse{
		SleepTime: bedtime.SleepTime,
		WakeTime:  bedtime.WakeTime,
	}

	return &bedtimeResponse, nil
}

func (s bedtimeService) UpdateBedtime(gId string, bedtime *BedtimeUpdater) error {
	err := s.bedtimeRepo.UpdateBedtime(gId, &repository.Bedtime{
		SleepTime: bedtime.SleepTime,
		WakeTime:  bedtime.WakeTime,
	})
	if err != nil {
		return err
	}
	return nil
}
