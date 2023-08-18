package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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
)

const (
	DBMaxOpenConnection     = 25
	DBMaxIdleConnection     = 25
	DBMaxConnectionLifeTime = 10 * time.Minute
)

type App struct {
	Conf   *config.Config
	Router *mux.Router
	DB     *sqlx.DB
}

func NewApp() *App {
	return &App{}
}

func Start() error {
	a := NewApp()
	c := config.NewConfig()
	c.Init()

	err := a.RunServer()
	if err != nil {
		logrus.Fatalf("")
		return err
	}

	a.MustPostgresConnection()
	st := store.NewStore(a.DB)
	srv := service.NewService(st, a.Conf)
	a.RegisterRouters(srv)

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
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	go func() {
		if err = server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server error: %v", err)
		}
	}()

	go func() {
		sig := <-cancelChan
		logrus.Printf("Caught signal %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
