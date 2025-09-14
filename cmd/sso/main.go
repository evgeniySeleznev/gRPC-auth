package main

import (
	"github.com/evgeniySeleznev/gRPC-auth/internal/app"
	"github.com/evgeniySeleznev/gRPC-auth/internal/config"
	"github.com/evgeniySeleznev/gRPC-auth/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// запуск с локальным конфигом
// go run cmd/sso/main.go --config=./config/local.yaml
func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting server", slog.Any("cfg", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	//TODO: инициализировать приложение (app)

	//TODO: запустить gRPC-сервер приложения (run app)

	//Gracefull Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signalStop := <-stop
	log.Info("stopping application", slog.String("signal", signalStop.String()))
	
	application.GRPCSrv.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
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
