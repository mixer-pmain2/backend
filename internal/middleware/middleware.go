package middleware

import (
	"fmt"
	"net/http"
	"pmain2/internal/config"
	"pmain2/pkg/cache"
	"pmain2/pkg/logger"
	"pmain2/pkg/utils/jwt"
	"strings"
	"time"
)

var (
	INFO, _  = logger.New("app", logger.INFO)
	ERROR, _ = logger.New("app", logger.ERROR)
	appCache = cache.CreateCache(time.Minute, time.Minute)
)

func CheckAuth(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//username, password, ok := r.BasicAuth()
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		jwtT := jwt.JWT(config.AppConfig.SecretKey)

		if jwtT.IsValid(reqToken) {
			h.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token, charset="UTF-8"`)
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
