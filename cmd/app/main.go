package main

import (
	"fmt"

	"github.com/godovasik/tanki_docker_sql/internal/fetcher"
	"github.com/godovasik/tanki_docker_sql/internal/models"
)

//"github.com/godovasik/tanki_docker_sql/internal/fetcher"

func main() {
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
	fmt.Println(data.Response.Name)

	var datastamp models.Datastamp

	datastamp.ConvertResponseToDatastamp(data)
	datastamp.NewPrint()

}
