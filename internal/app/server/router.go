package server

import (
	"net/http"

	"github.com/go-chi/chi"
	cmiddleware "github.com/go-chi/chi/middleware"
	phttp "github.com/valbury-repos/gotik/http"
	"github.com/gobliggg/berekbek/internal/app/handler"
)

// Router a chi mux
func Router(opt handler.HandlerOption, handlerCtx phttp.HandlerContext) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cmiddleware.RequestID)
	r.Use(cmiddleware.RealIP)
	r.Use(cmiddleware.Recoverer)

	// the handler
	phandler := phttp.NewHttpHandler(handlerCtx)

	healthCheckHandler := handler.HealthCheckHandler{}

	healthCheckHandler.HandlerOption = opt
	healthCheckHandler.Handler = phandler(healthCheckHandler.HealthCheck)

	// Setup your routing here
	r.Method(http.MethodGet, "/health-check", healthCheckHandler)

	r.Handle("/socket.io/", opt.SocketIO)

	return r
}