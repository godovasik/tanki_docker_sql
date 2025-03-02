package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/godovasik/tanki_docker_sql/internal/fetcher"
	"github.com/godovasik/tanki_docker_sql/internal/models"
	"github.com/godovasik/tanki_docker_sql/internal/storage"
	"github.com/godovasik/tanki_docker_sql/logger"
	"github.com/robfig/cron/v3"
)

type datastampAndUserId struct {
	datastamp *models.Datastamp
	user_id   int
}

func UpdateTask(ctx context.Context, userRepo storage.UserRepository, f fetcher.Fetcher) error {

	/*
	   1. Подтянуть откуда-то список игроков, из дб видимо
	   2. запустить горутину на каждого пользователя и загружать его данные в структурку
	   3. записать все это дело в дб
	*/

	users, err := userRepo.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	ch := make(chan datastampAndUserId, len(users))

	wg := sync.WaitGroup{}
	for _, user := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := f.SendRequest(ctx, user.Name) // надо обработать если юзер не найден
			if err != nil {
				logger.Log.Error("ошибка при отправке реквеста", err)
				return
			}
			// fmt.Println("resp:", resp)
			rawData, err := f.ParseResponse(resp)
			if err != nil {
				logger.Log.Error("ошибка парсинга: ", err)
				return
			}
			data := models.ConvertResponseToDatastamp(rawData)
			ch <- datastampAndUserId{data, user.Id}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for data := range ch {
		err = userRepo.UpdateDataForUser(ctx, data.datastamp, data.user_id)
		if err != nil {
			logger.Log.Error("ошибка при добавлении статы:", err)
		}
	}

	return nil

}

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

	ctx := context.Background()

	fmt.Println("yoooooooooooooooooo")

	UpdateTask(ctx, userRepo, f)

	logger.Log.Debug("initializing cron task")
	c := cron.New()
	_, err = c.AddFunc("0 */2 * * *", func() {
		_ = UpdateTask(ctx, userRepo, f)
	})
	if err != nil {
		logger.Log.Error(err)
	}

	c.Start()
	logger.Log.Info("task scheduled")
	select {}

}
