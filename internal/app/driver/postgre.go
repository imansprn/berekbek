package driver

import (
  "fmt"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"

  _ "github.com/lib/pq"
)

type DBPostgreOption struct {
  Host        string
  Port        int
  Username    string
  Password    string
  Name        string
  MaxPoolSize int
}

func NewPostgreDatabase(option DBPostgreOption) (gormDB *gorm.DB, err error) {
  gormDB, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", option.Host, option.Port, option.Username, option.Name, option.Password)), &gorm.Config{})
  if err != nil {
    return
  }

  return
}
