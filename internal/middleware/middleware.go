package middleware

import (
	"fmt"
	"net/http"
	"pmain2/internal/controller"
	"pmain2/pkg/cache"
	"pmain2/pkg/logger"
	"pmain2/pkg/utils"
	"time"
)

var (
	INFO, _  = logger.New("app", logger.INFO)
	ERROR, _ = logger.New("app", logger.ERROR)
	appCache = cache.CreateCache(time.Minute, time.Minute)
)

type auth struct {
	username string
	password string
}

func CheckAuth(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()
		if ok && payed {
			isAuth, ok := appCache.Get(auth{username: username, password: password})
			if !ok {
				c := controller.Init()
				var err error
				user, _ := utils.ToWin1251(username)
				pass, _ := utils.ToWin1251(password) // utils.ToASCII(password)
				isAuth, err = c.User.IsAuth(user, pass)
				if err != nil {
					INFO.Println("BasicAuth, ok=", isAuth, " err=", err)
					http.Error(w, err.Error(), http.StatusUnauthorized)
					ERROR.Println(err.Error())
					return
				}
				appCache.Set(auth{username: username, password: password}, isAuth, time.Second*10)
			}
			if isAuth != nil && isAuth.(bool) {
				h.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		INFO.Println(fmt.Sprint(r.URL))
		h.ServeHTTP(w, r)
	})
}

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func JsonHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
