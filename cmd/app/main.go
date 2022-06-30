package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pmain2/internal/application"
	"pmain2/internal/config"
	"pmain2/internal/database"
	"pmain2/internal/models"
	"pmain2/internal/server"
	"pmain2/pkg/cache"
)

func init() {
	application.InitLogger()
	var err error
	config.AppConfig, err = config.Create()
	if err != nil {
		application.ERROR.Println(err.Error())
	}

	conn, err := database.Connect()
	if err != nil {
		application.ERROR.Println(err.Error())
	}
	models.Model = models.Init(conn.DB)
	//defer conn.Close()

	cache.AppCache = cache.CreateCache(time.Minute, time.Minute)
}

func main() {
	application.INFO.Print("Star app\n")

	application.INFO.Println("Init handlers")
	server := server.Create(config.AppConfig)
	application.CreateRouters(server)

	go func() {
		err := server.Run()
		if err != nil {
			application.INFO.Println("Error started server: ", err)
			fmt.Printf("Error started server: %s \n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	application.INFO.Println("Stop app")
}
