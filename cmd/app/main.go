package main

import (
	"fmt"
	"net/http"

	"pmain2/internal/api"
	"pmain2/internal/apperror"
	"pmain2/internal/config"
	"pmain2/internal/middleware"
	"pmain2/internal/server"
	"pmain2/internal/web"
)

func main() {
	conf, err := config.Create()
	if err != nil {
		fmt.Println(err)
	}

	apiHandlers := api.Init()

	server := server.Create(conf)
	apiRouter := server.Router.PathPrefix("/api/v0").Subrouter()
	apiRouter.Use(middleware.CORS)
	apiRouter.Use(middleware.BasicAuth)
	apiRouter.Use(middleware.JsonHeader)
	apiRouter.Use(middleware.Logging)

	//User
	apiRouter.HandleFunc("/user/{id:[0-9]*}/prava/", apperror.Middleware(apiHandlers.User.GetPrava))
	apiRouter.HandleFunc("/user/{id:[0-9]*}/", apperror.Middleware(apiHandlers.User.GetUser))
	apiRouter.HandleFunc("/auth/signin/", apperror.Middleware(apiHandlers.User.Signin))
	//Patient
	apiRouter.HandleFunc("/patient/find/",
		apperror.Middleware(apiHandlers.Patient.Find))
	apiRouter.HandleFunc("/patient/{id:[0-9]}/", apperror.Middleware(apiHandlers.Patient.Get))
	apiRouter.HandleFunc("/spr/podr/", apperror.Middleware(apiHandlers.Spr.GetPodr))

	webRouter := server.Router.PathPrefix("/").Subrouter()
	webRouter.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./")))
	webRouter.Use(middleware.Logging)

	web.RoutesFrontend(webRouter, web.IndexServe)

	err = server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
