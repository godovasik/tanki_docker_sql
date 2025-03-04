package main

import (
	"context"
	"fmt"
	"time"

	"github.com/godovasik/tanki_docker_sql/internal/fetcher"
	"github.com/godovasik/tanki_docker_sql/internal/service"
	"github.com/godovasik/tanki_docker_sql/internal/storage"
	"github.com/godovasik/tanki_docker_sql/logger"
	"github.com/robfig/cron/v3"
)

func main() {
	//настройка логов
	logger.SetupLogger()

	logger.Log.Debug("logger initialized")

	//покдлючаемся к дб, на выходе - userRepo для тасок с дб и context
	userRepo, cleanup, err := storage.ConnectToDb()
	if err != nil {
		logger.Log.Error("cant conenct to db:", err)
		return
	}
	defer cleanup()

	f := fetcher.NewHTTPFetcher(10 * time.Second)

	srvc := service.NewUserService(userRepo, f)

	ctx := context.Background()

	fmt.Println("yoooooooooooooooooo")

	// она обрабатывает ошибки сама так что похуй я не знаю как правильно
	srvc.UpdateTask(ctx)

	logger.Log.Debug("initializing cron task")
	c := cron.New()
	_, err = c.AddFunc("0 */2 * * *", func() {
		_ = srvc.UpdateTask(ctx)
	})
	if err != nil {
		logger.Log.Error(err)
	}

	c.Start()
	logger.Log.Info("task scheduled")
	select {}

}
