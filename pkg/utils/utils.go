package utils

import (
	"io/ioutil"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func ToUTF8(text string) (string, error) {
	sr := strings.NewReader(text)
	tr := transform.NewReader(sr, charmap.Windows1251.NewDecoder())
	buf, err := ioutil.ReadAll(tr)
	if err != err {
		return "", err
	}

	text = string(buf) // строка в UTF-8
	return text, nil
}

func ToWin1251(text string) (string, error) {
	sr := strings.NewReader(text)
	tr := transform.NewReader(sr, charmap.Windows1251.NewEncoder())
	buf, err := ioutil.ReadAll(tr)
	if err != err {
		return "", err
	}

	text = string(buf) // строка в Win-1251
	return text, nil
}

func ToDate(t time.Time) string {
	return t.Format("2006-01-02")
}
