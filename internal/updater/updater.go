package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"path/filepath"
	"pmain2/internal/application"
	"pmain2/pkg/logger"
	"pmain2/pkg/version"
	"strings"
	"time"
)

var (
	l      *logger.Logger
	remote string
)

func init() {
	l, _ = logger.Create("updater")
}

const (
	updateFile   = "update.json"
	updateFolder = "update"
	oldFolder    = "old"
)

var (
	appStruct = []string{
		"migrations",
		"static",
		"frontend.routes",
		"index.html",
		"main.exe",
	}
)

type runUpdater interface {
	GetDuration() time.Duration
	IsEnabled() bool
	GetServiceName() string
	GetRemote() string
}

type fileUpdate struct {
	Version version.Version `json:"version"`
	Files   []string        `json:"files"`
}

func Run(t runUpdater) error {
	remote = t.GetRemote()
	go func() {
		for true {
			var client http.Client

			url := remote + "/" + updateFile
			fmt.Println(url)
			resp, err := client.Get(url)
			if err != nil {
				l.ERROR.Println(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				byteBody, err := io.ReadAll(resp.Body)
				if err != nil {
					l.ERROR.Println(err)
				}
				fileBody := fileUpdate{}
				err = json.Unmarshal(byteBody, &fileBody)
				if err != nil {
					l.ERROR.Println(err)
				}
				if ok, err := application.Version.IsHigh(fileBody.Version); !ok {
					l.INFO.Printf("update version is High, old %s, new %s", application.Version, fileBody.Version)
					l.INFO.Printf("uploading update files")
					err = uploadNewFiles(fileBody.Files)
					if err != nil {
						l.ERROR.Println(err)
					} else {
						err = saveCurrentVersion()
						if err != nil {
							l.ERROR.Println(err)
						} else {
							err = updateApp(fileBody.Files)
							if err != nil {
								l.ERROR.Println(err)
							}
						}
					}
					if err != nil {
						l.INFO.Printf("app updated")
					}
				} else if err != nil {
					l.ERROR.Println(err)
				}

			}
			<-time.Tick(t.GetDuration())
		}

	}()
	return nil
}

func uploadNewFiles(files []string) error {
	var client http.Client

	for _, f := range files {
		url, _ := url2.JoinPath(remote, f)
		fmt.Println(url)
		resp, err := client.Get(url)
		if err != nil {
			return err
		}
		byteBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		f, err := createFile(filepath.Join(updateFolder, f))
		if err != nil {
			log.Fatalln(err)
			return err
		}
		_, err = f.Write(byteBody)
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	return nil
}

func saveCurrentVersion() error {
	oldPath := filepath.Join(".", oldFolder, string(application.Version))
	for _, _f := range appStruct {
		path := filepath.Join(".", _f)
		if ok, err := isDir(path); ok {
			files, _ := filesInDir(path)
			for _, f := range files {
				_, err := copy(filepath.Join(oldPath, f), f)
				if err != nil {
					return err
				}
			}

		} else if err != nil {

		} else {
			_, err := copy(filepath.Join(oldPath, _f), path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func updateApp(files []string) error {
	for _, f := range files {
		fmt.Println(f)
	}

	return nil
}

func createFile(path string) (*os.File, error) {
	pathSplited := strings.Split(path, "\\")
	pathDir := strings.Join(pathSplited[0:len(pathSplited)-1], "\\")

	err := os.MkdirAll(pathDir, os.ModePerm)
	if err != nil {
		l.ERROR.Println(err)
		return nil, err
	}
	return os.Create(path)
}

func filesInDir(path string) ([]string, error) {
	var result []string
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			_files, err := filesInDir(filepath.Join(path, f.Name()))
			if err != nil {
				return nil, err
			}
			result = append(result, _files...)
		} else {
			result = append(result, filepath.Join(path, f.Name()))
		}
	}

	return result, nil
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func copy(currentPath, newPath string) (bool, error) {
	oldFile, _ := createFile(currentPath)
	currentFile, _ := os.Open(newPath)
	_, err := io.Copy(oldFile, currentFile)
	if err != nil {
		l.ERROR.Println(err)
		return false, err
	}
	currentFile.Close()
	oldFile.Close()

	return true, nil
}
