package main

import (
	"eda/models"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	models.ConnectDb()

	err := executeSQLScripts(models.DB, "../../sql")
	if err != nil {
		log.Fatalf("Ошибка выполнения SQL-скриптов: %v", err)
	}
}

func executeSQLScripts(db *gorm.DB, path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			result := db.Exec(string(content))
			if result.Error != nil {
				return result.Error
			}
			log.Printf("Успешно выполнен SQL-скрипт: %s", info.Name())
		}

		return nil
	})
}
