package main

import (
	"context"
	"fmt"

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
	userRepo, cleanup, err := storage.ConnectToDb()
	if err != nil {
		logger.Log.Error("cant conenct to db:", err)
		return
	}
	defer cleanup()
	ctx := context.Background()

	date, err := userRepo.FindLastStampDate(ctx, 1)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	fmt.Println(date)

	//про это пока забыли
	// data, err := userRepo.FindLastChangedDatastamp(ctx, 2)
	// if err != nil {
	// 	logger.Log.Error(err)
	// 	return
	// }
	// fmt.Print(data)

	// закончили подключаться и создали контекст.

	// вставка стампов
	// data := models.Datastamp{Name: "silly", Deaths: 17, Rank: 17, Kills: 34, EarnedCrystals: 1007}
	// err = userRepo.AddDatastamp(ctx, data, 2)
	// if err != nil {
	// 	logger.Log.Error(err)
	// 	return
	// }

	//getuserbyid
	// user, err := userRepo.GetUserById(ctx, 2)
	// if err != nil {
	// 	logger.Log.Error(err)
	// }
	// fmt.Print(user)

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
