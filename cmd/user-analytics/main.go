package main

import (
	"log/slog"
	"os"
	"user-analytics/internal/config"
	"user-analytics/internal/lib/logger/sl"
	"user-analytics/internal/lib/logger/slogpretty"
	"user-analytics/internal/storage/postgres"

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




	//TODO: init kafka


	//TODO: init chi router


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

