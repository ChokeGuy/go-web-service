package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"web-service/config"
	googledrive "web-service/pkg/google-drive"
	"web-service/pkg/handler"
	"web-service/pkg/kafka"
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

var wg sync.WaitGroup

func loadEnv() {
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
}

func setupRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// Middlewares
	r.Use(middlewares.CorsMiddleware)
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.ErrorHandlerMiddleware)
	r.Use(middlewares.RecoverMiddleware)

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

func waitForShutdown(ctx context.Context, srv *http.Server) {
	// Handle system signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-ctx.Done():
	case <-c:
		log.Println("Shutdown signal received")
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server Shutdown failed: %v", err)
	}
}

func initGoogleDrive() {
	googledrive.Init()
}

func initKafka(ctx context.Context) {
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := kafka.InitProducer(config.Env.KafkaBrokers); err != nil {
			log.Fatalf("Failed to initialize Kafka producer: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := kafka.InitConsumer(config.Env.KafkaBrokers, config.Env.KafkaGroupID); err != nil {
			log.Fatalf("Failed to initialize Kafka consumer: %v", err)
		}
	}()

	// Wait for Kafka to be ready or context cancel
	go func() {
		select {
		case <-ctx.Done():
			log.Println("Kafka initialization canceled")
		}
	}()
}

func main() {
	loadEnv()

	cfg := getServerConfig()
	router := setupRouter()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Google Drive and Kafka
	initGoogleDrive()
	initKafka(ctx)

	// Start HTTP server
	srv := startServer(cfg, router)

	// Wait for shutdown signal
	waitForShutdown(ctx, srv)

	// Wait for other services to close
	wg.Wait()
	log.Println("All services gracefully stopped")
}
