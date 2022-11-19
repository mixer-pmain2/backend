package config

import (
	"encoding/json"
	"io"
	"os"
	"pmain2/pkg/logger"
	"time"
)

const (
	filename = "settings.json"
)

var (
	l *logger.Logger
)

func init() {
	l, _ = logger.Create("settings")
}

type Updater struct {
	Enabled     bool   `json:"enabled"`
	Duration    int64  `json:"duration"`
	ServiceName string `json:"serviceName"`
	Remote      string `json:"remote"`
}

type AppSetting struct {
	Name string `json:"name"`
}

type Settings struct {
	App     AppSetting `json:"app"`
	Updater Updater    `json:"updater"`
}

func GetSettings() Settings {
	s := Settings{
		App: AppSetting{
			Name: "pmain2test",
		},
		Updater: Updater{
			Duration:    600,
			Enabled:     true,
			ServiceName: "pmain2",
			Remote:      "http://localhost:8091/pmain/update",
		},
	}
	jsonFile, err := os.Open(filename)
	if err != nil {
		l.ERROR.Println("Open settings.json failed, load default settings %s, error: %s", s, err)
		b, _ := json.Marshal(s)
		f, err := os.Create("./" + filename)
		if err != nil {
			l.ERROR.Println(err)
		}
		defer f.Close()
		f.Write(b)
	} else {
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			l.ERROR.Println(err)
		}

		err = json.Unmarshal(byteValue, &s)
		if err != nil {
			l.ERROR.Println(err)
		}
	}
	defer jsonFile.Close()

	return s
}

func (s Settings) GetUpdater() Updater {
	return s.Updater
}

func (s Settings) GetAppName() string {
	return s.App.Name
}

func (u Updater) GetDuration() time.Duration {
	return time.Second * time.Duration(u.Duration)
}

func (u Updater) GetServiceName() string {
	return u.ServiceName
}

func (u Updater) IsEnabled() bool {
	return u.Enabled
}

func (u Updater) GetRemote() string {
	return u.Remote
}
