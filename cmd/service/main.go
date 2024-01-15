package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"log/slog"

	"github.com/davecgh/go-spew/spew"
	"github.com/waiq/example-service/pkg/config"
	"github.com/waiq/example-service/pkg/controllers"
	"github.com/waiq/example-service/pkg/repository"
	"github.com/waiq/example-service/pkg/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	config.InitConfig()
	ctx := context.Background()

	database, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN:                  config.GetPostgresDSN(),
			PreferSimpleProtocol: true,
		},
	))

	if err != nil {
		panic(err)
	}

	repo, err := repository.New(ctx, database)
	if err != nil {
		panic(err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	if err := repo.Migration(ctx); err != nil {
		logger.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	c := controllers.NewBooksController(service.NewBookService(ctx, repo))
	svc := http.Server{
		Addr:    fmt.Sprintf(":%d", *config.ApplicationPort),
		Handler: c.Handler(),
	}

	spew.Dump(svc.Addr)

	go func() {
		defer wg.Done()

		logger.Info("http server starting")
		if err := svc.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("http server", err)
		}
	}()

	<-stopChan

	logger.Info("shutdown service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := svc.Shutdown(ctx); err != nil {
		logger.ErrorContext(ctx, "shutdown service", err)
	}

	wg.Wait()
}
