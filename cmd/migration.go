package cmd

import (
  "database/sql"
  "fmt"
  "os"
  "time"

  "github.com/golang-migrate/migrate/v4/database"
  "github.com/golang-migrate/migrate/v4/source"
  "github.com/golang-migrate/migrate/v4/source/file"
  "github.com/sirupsen/logrus"
  "github.com/gobliggg/berekbek/config"
  "gorm.io/gorm"

  "github.com/golang-migrate/migrate/v4"
  _ "github.com/golang-migrate/migrate/v4/source/file"
  "github.com/spf13/cobra"
  "github.com/gobliggg/berekbek/internal/app/appcontext"
)

const migrationDir = "migrations/sql"

var migrateCmd = &cobra.Command{
  Use: "migrate",
}

var migrateUpCmd = &cobra.Command{
  Use:   "up",
  Short: "Migrates the database to the most recent version available",
  Long:  `Migrates the database to the most recent version available`,
  RunE: func(cmd *cobra.Command, args []string) (err error) {
    c := config.Config()
    appCtx := appcontext.NewAppContext(c)

    if err = doMigrate(appCtx, c.GetString("database.connection"), source.Up, 0); err != nil {
      return
    }
    return
  },
}

var migrateDownCmd = &cobra.Command{
  Use:   "down",
  Short: "Undo a database migrations",
  Long:  `Undo a database migrations`,
  RunE: func(cmd *cobra.Command, args []string) (err error) {
    c := config.Config()
    appCtx := appcontext.NewAppContext(c)
    step, _ := cmd.Flags().GetInt("step")

    if err = doMigrate(appCtx, c.GetString("database.connection"), source.Down, step); err != nil {
      return
    }
    return
  },
}

var migrateNewCmd = &cobra.Command{
  Use:   "new",
  Short: "Create a new migrations",
  Long:  `Create a new migrations`,
  Args:  cobra.MinimumNArgs(1),
  RunE: func(cmd *cobra.Command, args []string) (err error) {
    if err = createMigrationFile(args[0], source.Up); err != nil {
      return
    }

    if err = createMigrationFile(args[0], source.Down); err != nil {
      return
    }
    return
  },
}

var migrateForceCmd = &cobra.Command{
  Use:   "force",
  Short: "Fix and force version",
  Long:  `Fix and force version`,
  RunE: func(cmd *cobra.Command, args []string) (err error) {
    c := config.Config()
    appCtx := appcontext.NewAppContext(c)
    version, _ := cmd.Flags().GetInt("version")

    if err = forceVersion(appCtx, c.GetString("database.connection"), version); err != nil {
      return
    }
    return
  },
}

func init() {
  migrateDownCmd.PersistentFlags().IntP("step","s", 1, "The number of migrations to be reverted")
  migrateForceCmd.PersistentFlags().IntP("version","v", 0, "Set version but don't run migration (ignores dirty state)")

  rootCmd.AddCommand(migrateCmd)
  migrateCmd.AddCommand(migrateDownCmd)
  migrateCmd.AddCommand(migrateNewCmd)
  migrateCmd.AddCommand(migrateUpCmd)
  migrateCmd.AddCommand(migrateForceCmd)
}

func doMigrate(appCtx *appcontext.AppContext, dbDialect string, direction source.Direction, step int) (err error) {
  var db *gorm.DB
  db, err = appCtx.GetDBInstance(dbDialect)
  if err != nil {
    logrus.Fatalf("Error connection to DB | %v", err)
    return
  }

  var sqlDB *sql.DB
  sqlDB, err = db.DB()
  if err != nil {
    logrus.Errorf("Fail migrations | %v", err)
    return
  }
  defer sqlDB.Close()

  var dbDriver database.Driver
  dbDriver, err = appCtx.GetDBDriver(dbDialect, sqlDB)
  if err != nil {
    logrus.Errorf("Fail migrations | %v", err)
    return
  }

  var mSource source.Driver
  mSource, err = (&file.File{}).Open(fmt.Sprintf("file://%s", migrationDir))
  if err != nil {
    logrus.Errorf("Opening file error | %v", err)
    return
  }

  m, err := migrate.NewWithInstance("file", mSource, db.Config.Name(), dbDriver)
  if err != nil {
    logrus.Errorf("Migrate error | %v", err)
    return
  }

  switch direction {
  case source.Up:
    err = m.Up()
    break
  case source.Down:
    err = m.Steps(step * -1)
    break
  }

  if err != nil {
    logrus.Errorf("Migrate error | %v", err)
    return
  }

  logrus.Infof("Migrate Success")
  return
}

func createMigrationFile(mName string, mDirection source.Direction) (err error) {
  filename := fmt.Sprintf("%s_%s.%s.sql", time.Now().Format("20060102150405"), mName, mDirection)
  filepath := fmt.Sprintf("%s/%s", migrationDir, filename)

  f, err := os.Create(filepath)
  if err != nil {
    logrus.Errorf("Error create migrations file | %v", err)
    return
  }
  defer f.Close()

  logrus.Infof("New migrations file has been created: %s)", filepath)
  return
}

func forceVersion(appCtx *appcontext.AppContext, dbDialect string, version int) (err error) {
  var db *gorm.DB
  db, err = appCtx.GetDBInstance(dbDialect)
  if err != nil {
    logrus.Fatalf("Error connection to DB | %v", err)
    return
  }

  var sqlDB *sql.DB
  sqlDB, err = db.DB()
  if err != nil {
    logrus.Errorf("Fail migrations | %v", err)
    return
  }
  defer sqlDB.Close()

  var dbDriver database.Driver
  dbDriver, err = appCtx.GetDBDriver(dbDialect, sqlDB)
  if err != nil {
    logrus.Errorf("Fail migrations | %v", err)
    return
  }

  var mSource source.Driver
  mSource, err = (&file.File{}).Open(fmt.Sprintf("file://%s", migrationDir))
  if err != nil {
    logrus.Errorf("Opening file error | %v", err)
    return
  }

  m, err := migrate.NewWithInstance("file", mSource, db.Config.Name(), dbDriver)
  if err != nil {
    logrus.Errorf("Migrate error | %v", err)
    return
  }

  if version == 0 {
    var uVersion uint
    var dirty bool
    uVersion, dirty, err = m.Version()
    if err != nil {
      logrus.Errorf("Migrate error | %v", err)
      return
    }

    if dirty == true {
      version = int(uVersion)
    }
  }

  err = m.Force(version)
  if err != nil {
    logrus.Errorf("Migrate error | %v", err)
    return
  }

  logrus.Infof("Migrate Success")
  return
}
