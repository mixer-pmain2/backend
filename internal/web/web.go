package web

import (
	"bufio"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

func IndexServe(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("template not found"))
		return
	}
	tmpl.Execute(w, nil)
}

func RoutesFrontend(router *mux.Router, handle http.HandlerFunc) {
	path, _ := os.Getwd()
	fullpath := filepath.Join(path, "frontend.routes")
	file, _ := os.Open(fullpath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		router.HandleFunc(url, handle).Methods(http.MethodGet)
	}
}
