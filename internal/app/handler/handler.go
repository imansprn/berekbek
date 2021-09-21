package handler

import (
	"github.com/gobliggg/berekbek/internal/app/commons"
	"github.com/gobliggg/berekbek/internal/app/service"
)

// HandlerOption option for handler, including all service
type HandlerOption struct {
	commons.Options
	*service.Services
}
