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
	path := filepath.Join("templates", "index.html")
	tmpl, _ := template.ParseFiles(path)
	tmpl.Execute(w, nil)
}

func RoutesFrontend(router *mux.Router, handle http.HandlerFunc) {
	path, _ := os.Getwd()
	file, _ := os.Open(path + "./.froutes")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		router.HandleFunc(url, handle).Methods(http.MethodGet)
	}
}
