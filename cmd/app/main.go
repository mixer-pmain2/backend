package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pmain2/internal/api"
	"pmain2/internal/apperror"
	"pmain2/internal/config"
	"pmain2/internal/database"
	"pmain2/internal/middleware"
	"pmain2/internal/models"
	"pmain2/internal/server"
	"pmain2/internal/web"
	"pmain2/pkg/cache"
	"pmain2/pkg/logger"
)

var (
	INFO, _  = logger.New("app", logger.INFO)
	ERROR, _ = logger.New("app", logger.ERROR)
)

func init() {
	var err error
	config.AppConfig, err = config.Create()
	if err != nil {
		ERROR.Println(err.Error())
	}

	conn, err := database.Connect()
	if err != nil {
		ERROR.Println(err.Error())
	}
	models.Model = models.Init(conn.DB)
	//defer conn.Close()

	cache.AppCache = cache.CreateCache(time.Minute, time.Minute)
}

func main() {
	apiHandlers := api.Init()

	server := server.Create(config.AppConfig)
	apiRouter := server.Router.PathPrefix("/api/v0").Subrouter()
	apiRouter.Use(middleware.CORS)
	apiRouter.Use(middleware.BasicAuth)
	apiRouter.Use(middleware.JsonHeader)
	apiRouter.Use(middleware.Logging)

	apiRouter.HandleFunc("/user/{id:[0-9]*}/uch/", apperror.Middleware(apiHandlers.User.GetUch)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/user/{id:[0-9]*}/prava/", apperror.Middleware(apiHandlers.User.GetPrava)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/user/{id:[0-9]*}/", apperror.Middleware(apiHandlers.User.GetUser)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/auth/signin/", apperror.Middleware(apiHandlers.User.Signin)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/find/", apperror.Middleware(apiHandlers.Patient.Find)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/hospital/", apperror.Middleware(apiHandlers.Patient.HistoryHospital)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/uchet/", apperror.Middleware(apiHandlers.Patient.FindUchet)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/visit/", apperror.Middleware(apiHandlers.Patient.HistoryVisits)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/visit/", apperror.Middleware(apiHandlers.Patient.NewVisit)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/", apperror.Middleware(apiHandlers.Patient.Get)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/prof/", apperror.Middleware(apiHandlers.Patient.NewProf)).Methods(http.MethodPost, http.MethodOptions)

	apiRouterNonAuth := server.Router.PathPrefix("/api/v0").Subrouter()
	apiRouterNonAuth.Use(middleware.CORS)
	apiRouterNonAuth.Use(middleware.JsonHeader)
	apiRouterNonAuth.Use(middleware.Logging)
	apiRouterNonAuth.HandleFunc("/spr/podr/", apperror.Middleware(apiHandlers.Spr.GetPodr)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/prava/", apperror.Middleware(apiHandlers.Spr.GetPrava)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/visit/", apperror.Middleware(apiHandlers.Spr.GetSprVisit)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/diag/", apperror.Middleware(apiHandlers.Spr.GetSprDiags)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/service/", apperror.Middleware(apiHandlers.Spr.GetParams)).Methods(http.MethodGet, http.MethodOptions)

	webRouter := server.Router.PathPrefix("/").Subrouter()
	webRouter.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./")))
	webRouter.Use(middleware.Logging)

	web.RoutesFrontend(webRouter, web.IndexServe)

	INFO.Printf("Starting server\n")

	go func() {
		err := server.Run()
		if err != nil {
			INFO.Println("Error started server: ", err)
			fmt.Printf("Error started server: %s \n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	INFO.Println("Stop server")
}
