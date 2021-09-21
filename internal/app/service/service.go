package service

import (
	"github.com/gobliggg/berekbek/internal/app/commons"
	"github.com/gobliggg/berekbek/internal/app/repository"
)

// Option anything any service object needed
type Option struct {
	commons.Options
	*repository.Repository
}

// Services all service object injected here
type Services struct {
	HealthCheck IHealthCheck
}
