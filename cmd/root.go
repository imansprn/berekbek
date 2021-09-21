package cmd

import (
  "fmt"
  "github.com/googollee/go-socket.io/engineio"
  "github.com/googollee/go-socket.io/engineio/transport"
  "github.com/googollee/go-socket.io/engineio/transport/polling"
  "github.com/googollee/go-socket.io/engineio/transport/websocket"
  "gorm.io/gorm"
  "net/http"
  "os"

  "github.com/gobliggg/berekbek/config"
  "github.com/gobliggg/berekbek/internal/app/appcontext"
  "github.com/gobliggg/berekbek/internal/app/commons"
  "github.com/gobliggg/berekbek/internal/app/repository"
  "github.com/gobliggg/berekbek/internal/app/server"
  "github.com/gobliggg/berekbek/internal/app/service"
  socketio "github.com/googollee/go-socket.io"
  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  phttp "github.com/valbury-repos/gotik/http"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "berekbek",
  Short: "A brief description of your application",
  Long: `A longer description that spans multiple lines and likely contains
			examples and usage of using your application.`,
  Run: func(cmd *cobra.Command, args []string) {
    start()
  },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize()
}

func start() {
  cfg := config.Config()
  logger := logrus.New()
  logger.SetFormatter(&logrus.JSONFormatter{})
  logger.SetReportCaller(true)

  app := appcontext.NewAppContext(cfg)
  var err error

  var db *gorm.DB
  if cfg.GetString("database.connection") != "dialect" {
    db, err = app.GetDBInstance(cfg.GetString("database.connection"))
    if err != nil {
      logrus.Fatalf("Failed to start, error connect to database | %v", err)
      return
    }
  }

  handlerCtx := phttp.NewContextHandler()
  commons.InjectErrors(&handlerCtx)

  socketIO := socketio.NewServer(&engineio.Options{
    Transports: []transport.Transport{
      &polling.Transport{
        CheckOrigin: func(r *http.Request) bool {
          return true
        },
      },
      &websocket.Transport{
        CheckOrigin: func(r *http.Request) bool {
          return true
        },
      },
    },
  })


  socketIO.OnConnect("/", func(s socketio.Conn) error {
    s.SetContext("")
    fmt.Println("connected:", s.ID())
    return nil
  })

  socketIO.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
    fmt.Println("notice:", msg)
    s.Emit("reply", "have "+msg)
  })

  socketIO.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
    s.SetContext(msg)
    return "recv " + msg
  })

  socketIO.OnEvent("/", "bye", func(s socketio.Conn) string {
    last := s.Context().(string)
    s.Emit("bye", last)
    s.Close()
    return last
  })

  socketIO.OnError("/", func(s socketio.Conn, e error) {
    fmt.Println("meet error:", e)
  })

  socketIO.OnDisconnect("/", func(s socketio.Conn, reason string) {
    fmt.Println("closed", reason)
  })

  go socketIO.Serve()
  defer socketIO.Close()

  opt := commons.Options{
    Config:   cfg,
    Database: db,
    Logger:   logger,
    SocketIO: socketIO,
  }

  repo := wiringRepository(repository.Option{
    Options: opt,
  })

  svc := wiringService(service.Option{
    Options:    opt,
    Repository: repo,
  })

  // run app
  svr := server.NewServer(opt, svc, handlerCtx)
  svr.StartApp()
}

func wiringRepository(repoOption repository.Option) *repository.Repository {
  // wiring up all your repos here
  repo := repository.Repository{}

  return &repo
}

func wiringService(serviceOption service.Option) *service.Services {
  // wiring up all services
  hc := service.NewHealthCheck(serviceOption)

  svc := service.Services{
    HealthCheck: hc,
  }

  return &svc
}
