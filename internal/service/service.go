package service

import (
	"context"
	"sync"

	"github.com/godovasik/tanki_docker_sql/internal/fetcher"
	"github.com/godovasik/tanki_docker_sql/internal/models"
	"github.com/godovasik/tanki_docker_sql/internal/storage"
	"github.com/godovasik/tanki_docker_sql/logger"
)

type UserService struct {
	userRepo storage.UserRepository
	fetcher  fetcher.Fetcher
}

// костыль но поебать
type datastampAndUserId struct {
	datastamp *models.Datastamp
	user_id   int
}

func NewUserService(userRepo storage.UserRepository, fetcher fetcher.Fetcher) *UserService {
	return &UserService{userRepo: userRepo, fetcher: fetcher}
}

func (s *UserService) UpdateTask(ctx context.Context) error {

	/*
	   1. Подтянуть откуда-то список игроков, из дб видимо
	   2. запустить горутину на каждого пользователя и загружать его данные в структурку
	   3. записать все это дело в дб
	*/

	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	ch := make(chan datastampAndUserId, len(users))

	wg := sync.WaitGroup{}
	for _, user := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := s.fetcher.SendRequest(ctx, user.Name) // надо обработать если юзер не найден
			if err != nil {
				logger.Log.Error("ошибка при отправке реквеста", err)
				return
			}
			// fmt.Println("resp:", resp)
			rawData, err := s.fetcher.ParseResponse(resp)
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
		err = s.userRepo.UpdateDataForUser(ctx, data.datastamp, data.user_id)
		if err != nil {
			logger.Log.Error("ошибка при добавлении статы:", err)
		}
	}

	return nil

}
