package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"password-manager/internal/config"
	"password-manager/internal/handler"
	"password-manager/internal/service"
	"password-manager/internal/store"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

const (
	DBMaxOpenConnection     = 25
	DBMaxIdleConnection     = 25
	DBMaxConnectionLifeTime = 10 * time.Minute

	DefaultTimeOut = 15 * time.Second
)

type App struct {
	Conf   *config.Config
	Router *mux.Router
	DB     *sqlx.DB
}

func NewApp(c *config.Config) *App {
	return &App{Conf: c}
}

func Start() error {
	c := config.NewConfig()
	c.Init()

	app := NewApp(c)

	if err := app.RunServer(); err != nil {
		logrus.Fatalf("")

		return err
	}

	app.MustPostgresConnection()
	st := store.NewStore(app.DB)
	srv := service.NewService(st, app.Conf)
	app.RegisterRouters(srv)

	return nil
}

func (a *App) RegisterRouters(s *service.Service) {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/api/ping", handler.Ping(s))
}

func (a *App) RunServer() error {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	idleConnsClosed := make(chan struct{})
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	listener, err := net.Listen("tcp", a.Conf.RunAddressValue)
	if err != nil {
		logrus.Fatalf("failed to listen on address %v: %s", a.Conf.RunAddressValue, err.Error())
	}

	server := &http.Server{
		Handler:           a.Router,
		WriteTimeout:      DefaultTimeOut,
		ReadTimeout:       DefaultTimeOut,
		ReadHeaderTimeout: DefaultTimeOut,
		TLSConfig:         nil,
	}

	go func() {
		if err = server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server error: %v", err)
		}
	}()

	go func() {
		sig := <-cancelChan
		logrus.Printf("Caught signal %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeOut)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logrus.Fatalf("Server shutdown error: %v", err)

			return
		}
		close(idleConnsClosed)
	}()

	<-idleConnsClosed

	logrus.Info("Server shutdown successfully")
	return nil
}

func (a *App) MustPostgresConnection() {
	db, err := sqlx.Open("postgres", a.Conf.DatabaseURIValue)
	if err != nil {
		panic(err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		defer db.Close()
	}

	db.SetMaxOpenConns(DBMaxOpenConnection)
	db.SetMaxIdleConns(DBMaxIdleConnection)
	db.SetConnMaxLifetime(DBMaxConnectionLifeTime)

	a.DB = db
}
