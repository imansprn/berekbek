package server

import (
	"context"
	"fmt"
	phttp "github.com/valbury-repos/gotik/http"
	"net/http"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/gobliggg/berekbek/internal/app/commons"
	"github.com/gobliggg/berekbek/internal/app/handler"
	"github.com/gobliggg/berekbek/internal/app/service"
)

// IServer interface for server
type IServer interface {
	StartApp()
}

type server struct {
	opt        commons.Options
	services   *service.Services
	handlerCtx phttp.HandlerContext
}

// NewServer create object server
func NewServer(opt commons.Options, services *service.Services, handlerCtx phttp.HandlerContext) IServer {
	return &server{
		opt:        opt,
		services:   services,
		handlerCtx: handlerCtx,
	}
}

func (s *server) StartApp() {
	var srv http.Server
	idleConnectionClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		logrus.Infoln("[API] Server is shutting down")

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logrus.Infof("[API] Fail to shutting down: %v", err)
		}
		close(idleConnectionClosed)
	}()

	srv.Addr = fmt.Sprintf("%s:%d", s.opt.Config.GetString("app.host"), s.opt.Config.GetInt("app.port"))
	hOpt := handler.HandlerOption{
		Options:  s.opt,
		Services: s.services,
	}

	srv.Handler = Router(hOpt, s.handlerCtx)

	logrus.Infof("[API] HTTP serve at %s\n", srv.Addr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logrus.Infof("[API] Fail to start listen and server: %v", err)
	}

	<-idleConnectionClosed
	logrus.Infoln("[API] Bye")
}