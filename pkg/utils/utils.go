package utils

import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
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

func RuneToAscii(r rune) string {
	if r < 128 {
		return string(r)
	} else {
		return "\\u" + strconv.FormatInt(int64(r), 16)
	}
}

func ToASCII(text string) string {
	var res string
	for _, s := range text {
		res = res + RuneToAscii(s)
	}
	return res
}
