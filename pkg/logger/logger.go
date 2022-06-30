package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	INFO    = "INFO"
	ERROR   = "ERROR"
	WARNING = "WARNING"
)

type logWrite struct {
	file *os.File
}

func (l *logWrite) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return l.file.Write(p)
}

func New(name, lvl string) (*log.Logger, error) {
	path, err := os.Getwd()
	pathLogs := filepath.Join(path, "logs")
	err = os.MkdirAll(pathLogs, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	date := time.Now().Format("2006-01-02")
	filename := filepath.Join(pathLogs, name+"_"+date+".log")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return log.New(&logWrite{file: file}, fmt.Sprintf(`%s %s `, name, lvl), log.Ldate|log.Ltime|log.Lshortfile), nil
}
