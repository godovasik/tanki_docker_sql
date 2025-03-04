package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/godovasik/tanki_docker_sql/internal/fetcher"
	"github.com/godovasik/tanki_docker_sql/internal/models"
	"github.com/godovasik/tanki_docker_sql/internal/storage"
	"github.com/godovasik/tanki_docker_sql/logger"
	"github.com/jackc/pgx/v5"
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
			rawData, err := s.fetcher.GetUserStats(ctx, user.Name)
			if err == fmt.Errorf("NOT_FOUND") {
				logger.Log.Errorf("user %v not found, should delete", user.Name)
			}
			lastScore, err := s.userRepo.GetLastDatastampScore(ctx, user.Id)
			if err != nil {
				if err == pgx.ErrNoRows {
					logger.Log.Error("первый стамп:", user.Name)
				} else {
					logger.Log.Error("да пошел я нахуй")
				}
			}
			if lastScore != rawData.Response.Score {
				data := models.ConvertResponseToDatastamp(rawData)
				ch <- datastampAndUserId{data, user.Id}
			} else {
				logger.Log.Debug(user.Name, ": nothing to update")
			}
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

func (s *UserService) AddUser(ctx context.Context, username string) error {
	_, err := s.fetcher.GetUserStats(ctx, username)
	if err == fmt.Errorf("NOT_FOUND") {
		return fmt.Errorf("user %v not found", username)
	}
	err = s.userRepo.CreateUser(ctx, models.User{Name: username})
	return err
}
