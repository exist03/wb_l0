package service

import (
	"fmt"
	"strconv"
	"wb_l0/common"
	"wb_l0/pkg/logger"
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
	log := logger.GetLogger()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, common.ErrInvalidID
	}
	log.Debug().Msg(fmt.Sprintf("id == %d", idInt))
	bytes, err := s.repository.Get(idInt)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
