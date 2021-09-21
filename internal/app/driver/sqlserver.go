package driver

import (
  "fmt"
  "gorm.io/driver/sqlserver"
  "gorm.io/gorm"

  _ "github.com/lib/pq"
)

type DBSQLServerOption struct {
  Host     string
  Port     int
  Username string
  Password string
  Name     string
}

func NewSQLServerDatabase(option DBSQLServerOption) (gormDB *gorm.DB, err error) {
  gormDB, err = gorm.Open(sqlserver.Open(fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", option.Username, option.Password, option.Host, option.Port, option.Name)), &gorm.Config{})
  if err != nil {
    return
  }

  return
}
