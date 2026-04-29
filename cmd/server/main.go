package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tool-service/internal/config"
	httpapi "tool-service/internal/http"
	"tool-service/internal/repository"
	phonerepo "tool-service/internal/repository/phone"
	toollogrepo "tool-service/internal/repository/toollog"
	currencysvc "tool-service/internal/service/currency"
	phonesvc "tool-service/internal/service/phone"
	regionsvc "tool-service/internal/service/region"
	toollogsvc "tool-service/internal/service/toollog"
	wordstatsvc "tool-service/internal/service/wordstat"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.Load()
	if err != nil {
		log.Error("config", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := repository.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Error("db connect", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := repository.Migrate(ctx, db); err != nil {
		log.Error("migrate", "err", err)
		os.Exit(1)
	}

	tlRepo := toollogrepo.New(db)
	tlSvc := toollogsvc.New(tlRepo, log)

	phoneRepo := phonerepo.New(db)
	phoneSvc := phonesvc.New(phoneRepo, tlSvc)

	h := &httpapi.Handlers{
		WordStatSvc: wordstatsvc.New(cfg, tlSvc),
		PhoneSvc:    phoneSvc,
		CurrencySvc: currencysvc.New(cfg, tlSvc),
		RegionSvc:   regionsvc.New(tlSvc),
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpapi.NewRouter(h),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
