package service

import (
	"strconv"
	"wb_l0/common"
)

type repository interface {
	Get(id int) ([]byte, error)
}
type Service struct {
	repository
}

func New(repository repository) *Service {
	return &Service{repository}
}

func (s *Service) Get(id string) ([]byte, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, common.ErrInvalidID
	}
	bytes, err := s.repository.Get(idInt)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
