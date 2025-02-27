package main

import (
	"context"

	"github.com/godovasik/tanki_docker_sql/internal/storage"
	"github.com/godovasik/tanki_docker_sql/logger"
	// "github.com/robfig/cron/v3"
)

//"github.com/godovasik/tanki_docker_sql/internal/fetcher"

func main() {
	//настройка логов
	logger.SetupLogger()

	logger.Log.Debug("logger initialized")

	//покдлючаемся к дб, на выходе - userRepo для тасок с дб и context
	DSN := "postgres://tanki_enjoyer:r@172.24.125.42:5432/game_stats"
	logger.Log.Info("connecting to db...")

	cfg := storage.Config{DSN: DSN}
	pool, err := storage.NewPostgresPool(cfg)
	if err != nil {
		logger.Log.Error("db connection err:", err)
	}
	defer pool.Close()

	logger.Log.Info("we connected to db!")

	userRepo := storage.NewUserRepository(pool)
	ctx := context.Background()
	// закончили подключаться и создали контекст. надо бы все это обернуть

	// getAllUsers
	// users, err := userRepo.GetAllUsers(ctx)
	// if err != nil {
	// 	logger.Log.Error(err)
	// 	return
	// }
	// fmt.Println(users)

	// тестировал добавление юзера
	// err = userRepo.CreateUser(ctx, models.User{Name: "silly"})
	// if err != nil {
	// 	logger.Log.Error(err)
	// }
	// logger.Log.Info("we added some silly boi!")

	//ставим крон таску
	// logger.Log.Debug("initializing cron task")
	// c := cron.New()
	// _, err = c.AddFunc("0 3 * * *", UpdateTask)
	// if err != nil {
	// 	logger.Log.Error(err)
	// }
	// logger.Log.Info("task scheduled")

	//тестировал парсер
	/*
		resp, err := fetcher.SendRequest("silly")
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		// fmt.Println("resp:", resp)
		data, err := fetcher.ParseResponse(resp)
		if err != nil {
			fmt.Println(err)
			return
		}
		var datastamp models.Datastamp

		datastamp.ConvertResponseToDatastamp(data)
		datastamp.NewPrint(3)
	*/
	logger.Log.Infof("disconnected from db")
}

func UpdateTask() {
	/*
	   1. Подтянуть откуда-то список игроков, из дб видимо
	   2. запустить горутину на каждого пользователя и загружать его данные в структурку
	   3. записать все это дело в дб
	*/

}
