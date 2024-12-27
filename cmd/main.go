package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"web-service/config"
	"web-service/pkg/database"
	googledrive "web-service/pkg/google-drive"
	"web-service/pkg/handler"
	"web-service/pkg/middlewares"

	"github.com/gorilla/mux"
)

type ServerConfig struct {
	Host          string
	Port          string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
	ShutdownDelay time.Duration
}

func loadEnv() error {
	if err := config.Load(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func setupRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// Middlewares
	r.Use(middlewares.CorsMiddleware)
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.ErrorHandlerMiddleware)

	// Not Found and Not Allow Handler
	handler.NotFoundHandler(r)
	handler.NotAllowHandler(r)

	// Api V1
	apiV1Router := r.PathPrefix("/api/v1").Subrouter()
	handler.GoogleDriveRoutes(apiV1Router)
	handler.HomeRoutes(apiV1Router)
	handler.ProductRoutes(apiV1Router)

	return r
}

func getServerConfig() *ServerConfig {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections")
	flag.Parse()

	return &ServerConfig{
		Host:          config.Env.Host,
		Port:          config.Env.Port,
		ReadTimeout:   time.Second * 15,
		WriteTimeout:  time.Second * 15,
		IdleTimeout:   time.Second * 60,
		ShutdownDelay: wait,
	}
}

func startServer(cfg *ServerConfig, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv
}

func waitForShutdown(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server gracefully stopped")
}

func main() {
	loadEnv()

	cfg := getServerConfig()
	router := setupRouter()
	srv := startServer(cfg, router)

	client, err := database.MongoDBClient()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.DisconnectMongoDB(client)

	googledrive.Init()

	waitForShutdown(srv, cfg.ShutdownDelay)
}
