package middleware

import (
	"fmt"
	"log"
	"net/http"
	"pmain2/internal/api"
	"pmain2/internal/controller"
	"pmain2/pkg/logger"
)

var (
	INFO, _  = logger.New("app", logger.INFO)
	ERROR, _ = logger.New("app", logger.ERROR)
)

func BasicAuth(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()
		if ok {
			c := controller.Init()
			isAuth, err := c.User.IsAuth(username, password)
			INFO.Println("BasicAuth, ok=", isAuth, " err=", err)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				ERROR.Println(err.Error())
				return
			}
			if isAuth {
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
		api.INFO.Println(fmt.Sprint(r.URL))
		log.Println(r.URL)
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
