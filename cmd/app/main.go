package main

import (
	"fmt"
	"os"
	"os/signal"
	"pmain2/internal/migration"
	"pmain2/internal/report"
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
	application.INFO.Println("Init apps...")
	var err error
	application.INFO.Println("Load config")
	config.AppConfig, err = config.Create()
	if err != nil {
		application.ERROR.Println(err.Error())
	}

	conn, err := database.Connect()
	if err != nil {
		application.ERROR.Println(err.Error())
	}
	models.Init(conn.DB)
	err, tx := models.Model.CreateTx()
	if err != nil {
		application.ERROR.Println(err)
	}
	defer tx.Rollback()
	migration.Init(tx)
	err, tx = models.Model.CreateTx()
	if err != nil {
		application.ERROR.Println(err)
	}
	defer tx.Rollback()
	migration.LoadMigrations(tx)

	cache.AppCache = cache.CreateCache(time.Minute, time.Minute)
}

func main() {
	application.INFO.Print("Star app\n")

	application.INFO.Println("Init handlers")
	server := server.Create(config.AppConfig)
	application.CreateRouters(server)

	go func() {
		application.INFO.Println("Start http server")
		err := server.Run()
		if err != nil {
			application.INFO.Println("Error started server: ", err)
			fmt.Printf("Error started server: %s \n", err)
		}
	}()
	go func() {
		application.INFO.Println("Start report runner")
		err := report.Run()
		if err != nil {
			application.INFO.Println("Error started report runner: ", err)
			fmt.Printf("Error started report runner: %s \n", err)
		}
	}()
	//go func() {
	//	settings := config.GetSettings()
	//	application.INFO.Println("Start updater")
	//	err := updater.Run(settings.GetUpdater())
	//	if err != nil {
	//		application.INFO.Println("Error started updater: ", err)
	//		fmt.Printf("Error started updater: %s \n", err)
	//	}
	//}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	application.INFO.Println("Stop app")
}
