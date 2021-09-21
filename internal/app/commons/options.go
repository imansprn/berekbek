package commons

import (
  "github.com/gobliggg/berekbek/config"
  socketio "github.com/googollee/go-socket.io"
  "github.com/sirupsen/logrus"
  "gorm.io/gorm"
)

// Options common option for all object that needed
type Options struct {
  Config   config.Provider
  Database *gorm.DB
  Logger   *logrus.Logger
  SocketIO *socketio.Server
}
