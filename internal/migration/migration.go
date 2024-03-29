package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"pmain2/internal/application"
	"sort"
	"strconv"
	"strings"
)

const (
	migrationsPath = "migrations"
)

func Init(tx *sql.Tx) {
	application.INFO.Println("Init migration")
	sqlQuery := fmt.Sprintf("CREATE TABLE MIGRATIONS (ID INTEGER GENERATED BY DEFAULT AS IDENTITY NOT NULL,NAME VARCHAR(200) NOT NULL,CONSTRAINT MIGRATIONS_PK PRIMARY KEY (ID))")
	_, err := tx.Exec(sqlQuery)
	defer tx.Rollback()
	if err != nil {
		application.ERROR.Println(err)
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		application.ERROR.Println(err)
		return
	}

}

func getLastMigration(tx *sql.Tx) (int, error) {
	sqlQuery := fmt.Sprintf("select first 1 name from migrations order by id desc")
	row := tx.QueryRow(sqlQuery)
	name := ""
	err := row.Scan(&name)
	if err != nil {
		application.ERROR.Println(err)
		return 0, err
	}
	id, _ := indexFromName(name)
	return id, nil
}

func indexFromName(name string) (int, error) {
	nameSplit := strings.Split(strings.Split(name, ".")[0], "_")
	id, err := strconv.Atoi(nameSplit[0])
	if err != nil {
		application.ERROR.Println(err)
		return 0, err
	}
	return id, nil
}

func LoadMigrations(tx *sql.Tx) {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		application.ERROR.Println(err)
	}

	sort.Slice(files, func(i, j int) bool {
		a, _ := indexFromName(files[i].Name())
		b, _ := indexFromName(files[j].Name())
		return a < b
	})

	lIndex, _ := getLastMigration(tx)
	for _, file := range files {
		if isSql := strings.HasSuffix(file.Name(), ".sql"); !isSql {
			continue
		}
		index, _ := indexFromName(file.Name())
		if index > lIndex {
			lIndex = index
			pathFile := path.Join(migrationsPath, file.Name())
			f, err := os.Open(pathFile)
			if err != nil {
				application.ERROR.Println(err)
				tx.Rollback()
				break
			}
			data, err := ioutil.ReadAll(f)
			f.Close()
			if err != nil {
				application.ERROR.Println(err)
				tx.Rollback()
				break
			}
			_, err = tx.Exec(string(data))
			if err != nil {
				application.ERROR.Println(err)
				tx.Rollback()
				break
			}
			application.INFO.Println("load migration", file.Name(), "\n\t", string(data))
			_, err = tx.Exec(fmt.Sprintf("insert into migrations (name) values ('%s')", strings.Split(file.Name(), ".")[0]))
			if err != nil {
				application.ERROR.Println(err)
				tx.Rollback()
				break
			}
		}
	}
	tx.Commit()
}
