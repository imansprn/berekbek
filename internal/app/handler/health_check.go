package handler

import (
	"github.com/valbury-repos/gotik/structs"
	"net/http"
)

// HealthCheckHandler object for health check handler
type HealthCheckHandler struct {
	HandlerOption
	http.Handler
}

// HealthCheck checking if all work well
func (h HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) (data interface{}, meta structs.Meta, err error) {
	return
}
