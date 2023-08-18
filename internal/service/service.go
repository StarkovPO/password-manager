package service

import "password-manager/internal/config"

type StoreInterface interface {
}

type Service struct {
	store  StoreInterface
	config config.Config
}

func NewService(s StoreInterface, config *config.Config) *Service {
	return &Service{
		store:  s,
		config: *config,
	}
}
