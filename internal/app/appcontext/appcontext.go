package appcontext

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"
	"github.com/gobliggg/berekbek/config"
	"github.com/gobliggg/berekbek/internal/app/driver"
	"gorm.io/gorm"
)

const (
	// DBDialectMysql rdbms dialect name for MySQL
	DBDialectMysql = "mysql"

	// DBDialectPostgres rdbms dialect name for PostgreSQL
	DBDialectPostgres = "postgres"

	// DBDialectPostgres rdbms dialect name for SQLServer
	DBDialectSqlServer = "sqlserver"
)

// AppContext the app context struct
type AppContext struct {
	config config.Provider
}

// NewAppContext initiate appcontext object
func NewAppContext(config config.Provider) *AppContext {
	return &AppContext{
		config: config,
	}
}

func (a *AppContext) GetDBInstance(dbType string) (gormDB *gorm.DB, err error) {
	switch dbType {
	case DBDialectMysql:
		dbOption := a.getMysqlOption()
		gormDB, err = driver.NewMysqlDatabase(dbOption)
	case DBDialectPostgres:
		dbOption := a.getPostgreOption()
		gormDB, err = driver.NewPostgreDatabase(dbOption)
	case DBDialectSqlServer:
		dbOption := a.getSqlServerOption()
		gormDB, err = driver.NewSQLServerDatabase(dbOption)
	default:
		err = errors.New("error get db instance, unknown db type")
	}

	return
}

func (a *AppContext) GetDBDriver(dbType string, sqlDB *sql.DB) (dbDriver database.Driver, err error) {
	switch dbType {
	case DBDialectMysql:
		dbDriver, err = mysql.WithInstance(sqlDB, &mysql.Config{})
	case DBDialectPostgres:
		dbDriver, err = postgres.WithInstance(sqlDB, &postgres.Config{})
	case DBDialectSqlServer:
		dbDriver, err = sqlserver.WithInstance(sqlDB, &sqlserver.Config{})
	default:
		err = errors.New("error get db instance, unknown db type")
	}

	return
}

func (a *AppContext) getMysqlOption() driver.DBMysqlOption {
	return driver.DBMysqlOption{
		Host:     a.config.GetString("database.host"),
		Port:     a.config.GetInt("database.port"),
		Username: a.config.GetString("database.username"),
		Password: a.config.GetString("database.password"),
		Name:     a.config.GetString("database.name"),
	}
}

func (a *AppContext) getPostgreOption() driver.DBPostgreOption {
	return driver.DBPostgreOption{
		Host:     a.config.GetString("database.host"),
		Port:     a.config.GetInt("database.port"),
		Username: a.config.GetString("database.username"),
		Password: a.config.GetString("database.password"),
		Name:     a.config.GetString("database.name"),
	}
}

func (a *AppContext) getSqlServerOption() driver.DBSQLServerOption {
	return driver.DBSQLServerOption{
		Host:     a.config.GetString("database.host"),
		Port:     a.config.GetInt("database.port"),
		Username: a.config.GetString("database.username"),
		Password: a.config.GetString("database.password"),
		Name:     a.config.GetString("database.name"),
	}
}