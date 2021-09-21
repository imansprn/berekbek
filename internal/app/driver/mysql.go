package driver

import (
  "fmt"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "time"

  _ "github.com/go-sql-driver/mysql"
)

// DBMysqlOption options for mysql connection
type DBMysqlOption struct {
  Host                 string
  Port                 int
  Username             string
  Password             string
  Name                 string
  AdditionalParameters string
  MaxOpenConns         int
  MaxIdleConns         int
  ConnMaxLifetime      time.Duration
}

func NewMysqlDatabase(option DBMysqlOption) (gormDB *gorm.DB, err error) {
  gormDB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", option.Username, option.Password, option.Host, option.Port, option.Name, option.AdditionalParameters)), &gorm.Config{})
  if err != nil {
    return
  }
  return
}
