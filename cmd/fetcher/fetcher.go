package main

import (
	"context"
	"sync"
	"time"

	"github.com/godovasik/tanki_docker_sql/internal/fetcher"
	"github.com/godovasik/tanki_docker_sql/internal/models"
	"github.com/godovasik/tanki_docker_sql/internal/storage"
	"github.com/godovasik/tanki_docker_sql/logger"
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

	UpdateTask(ctx, userRepo, f)

	// logger.Log.Debug("initializing cron task")
	// c := cron.New()
	// _, err = c.AddFunc("0 */2 * * *", func() {
	// 	_ = UpdateTask(ctx, userRepo, f)
	// })
	// if err != nil {
	// 	logger.Log.Error(err)
	// }

	// c.Start()
	// logger.Log.Info("task scheduled")
	// select {}

	// fmt.Println(datastamp.Hulls["Wasp"])

	// userData, err := userRepo.FindLastGearStats(ctx, 2, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(userData)

	// err = userRepo.AddGearStats(ctx, 6, 4, models.GearData{TimePlayed: 15, ScoreEarned: 65})
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// find last
	// date, err := userRepo.FindLastStampDate(ctx, 2)
	// if err != nil {
	// 	logger.Log.Error(err)
	// 	return
	// }
	// fmt.Println(date)

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
}
