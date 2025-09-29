package main

import (
	"log/slog"
	"net/http"
	"os"
	"user-analytics/internal/config"
	"user-analytics/internal/lib/logger/sl"
	"user-analytics/internal/lib/logger/slogpretty"
	"user-analytics/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

// go run ./cmd/user-analytics/main.go

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)


func main() {
	
	initConfig()
	cfg := config.MustLoad()
	

	log := setupLogger(cfg.Env)

	log.Info(
		"starting user-analytics",
		slog.String("env", cfg.Env),
		slog.String("version", "666"),
	)
	log.Debug("debug messages enabled")

	db, err := postgres.New(cfg.DBConfig)

	if err!=nil{
		log.Error("failed to init storage", sl.Err(err))
	}
	
	_ = db
	
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

	log.Info("starting server", slog.String("adress",cfg.Adress))

	srv := &http.Server{
		Addr:	 cfg.Adress,
		Handler: router,
		ReadTimeout: 	cfg.HTTPServer.Timeout,
		WriteTimeout: 	cfg.HTTPServer.Timeout,
		IdleTimeout: 	cfg.HTTPServer.IdleTimeout,
	}
	
	if err:=srv.ListenAndServe();err!=nil{
		log.Error("failed to start server", sl.Err(err))
	}
	
	log.Error("server stopped")

	//TODO: init kafka

	//TODO: run server


}
func initConfig(){
	if err := godotenv.Load();err!=nil{
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

