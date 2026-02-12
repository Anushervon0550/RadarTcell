package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/Anushervon0550/RadarTcell/internal/config"
	"github.com/Anushervon0550/RadarTcell/internal/db"
	"github.com/Anushervon0550/RadarTcell/internal/httpapi"
)

func main() {
	_ = godotenv.Overload()
	// локально удобно, в проде можно не использовать

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load()
	if err != nil {
		log.Error("config error", "err", err)
		os.Exit(1)
	}
	//для завершения сервера по сигналу делает CNTRL+C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// DB
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("db connect error", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	// HTTP
	router := httpapi.NewRouter(httpapi.Deps{DB: pool})
	srv := httpapi.NewServer(cfg.AppPort, router)

	// стартуем сервер в отдельной горутине
	go func() {
		if err := srv.Start(log); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", "err", err)
			stop()
		}
	}()

	// ждём сигнал остановки
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(shutdownCtx, log)
	log.Info("bye")
}
