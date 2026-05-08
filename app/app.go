package app

import (
	"company-site/config"
	"company-site/handlers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Основной роутер
	mux := http.NewServeMux()
	handlers.SetupRoutes(mux, cfg)

	// Сервер с graceful shutdown
	addr := "0.0.0.0:" + cfg.ServerPort
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Printf("Сервер запущен на http://localhost:%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидание сигнала
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("Завершение работы сервера...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}
	log.Println("Сервер остановлен")

	return nil
}
