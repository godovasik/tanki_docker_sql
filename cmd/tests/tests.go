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