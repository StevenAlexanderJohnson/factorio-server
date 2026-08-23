package services

import "factorio/internal/config"

type SaveService struct {
	cfg *config.FactorioConfig
}

func NewSaveService(cfg *config.FactorioConfig) *SaveService {
	return &SaveService{
		cfg: cfg,
	}
}
