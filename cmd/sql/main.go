package main

import (
	"crypto/sha256"
	"eda/logger"
	"eda/models"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := initializeServices(); err != nil {
		logger.Log.Fatal("Failed to initialize services: ", zap.Error(err))
	}

	//err := executeSQLMigrations(models.DB, "../../sql/migrations")
	//if err != nil {
	//	logger.Log.Fatal("Failed to execute sql: ", zap.Error(err))
	//}

	err := executeSQLTriggers(models.DB, "../../sql/triggers")
	if err != nil {
		logger.Log.Fatal("Failed to execute sql trigger: ", zap.Error(err))
	}
}

func initializeServices() error {
	if err := logger.InitializeZapCustomLogger(); err != nil {
		return err
	}

	if err := models.ConnectDb(); err != nil {
		return err
	}

	return nil
}

//func executeSQLMigrations(db *gorm.DB, path string) error {
//	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//
//		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql") {
//			file, err := os.Open(path)
//			if err != nil {
//				return err
//			}
//			defer file.Close()
//
//			content, err := io.ReadAll(file)
//			if err != nil {
//				return err
//			}
//
//			result := db.Exec(string(content))
//			if result.Error != nil {
//				return result.Error
//			}
//			log.Printf("Успешно выполнен SQL-скрипт: %s", info.Name())
//		}
//
//		return nil
//	})
//}

func executeSQLTriggers(db *gorm.DB, path string) error {
	// Функция для получения хеша содержимого файла
	getHash := func(content []byte) string {
		hasher := sha256.New()
		hasher.Write(content)
		return hex.EncodeToString(hasher.Sum(nil))
	}

	// Функция для проверки и сохранения хеша скрипта
	checkAndSaveHash := func(db *gorm.DB, scriptName, scriptHash string) bool {
		var scriptExecution models.ScriptExecutions
		// Проверяем, есть ли запись о скрипте в базе данных
		result := db.Model(&models.ScriptExecutions{}).Where("script_name = ?", scriptName).First(&scriptExecution)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Ошибка при получении хеша из базы данных: %s", result.Error)
			return false
		}

		// Если хеш изменился или скрипт выполняется впервые
		if scriptExecution.ScriptHash == "" || scriptExecution.ScriptHash != scriptHash {
			// Обновляем хеш в базе данных
			if scriptExecution.ScriptHash == "" {
				db.Exec("INSERT INTO script_executions (script_name, script_hash) VALUES (?, ?)", scriptName, scriptHash)
			} else {
				db.Exec("UPDATE script_executions SET script_hash = ? WHERE script_name = ?", scriptHash, scriptName)
			}
			return true
		}
		return false
	}

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

			// Получаем хеш содержимого файла
			scriptHash := getHash(content)

			// Проверяем и сохраняем хеш скрипта
			if checkAndSaveHash(db, info.Name(), scriptHash) {
				// Хеш изменился, выполняем скрипт
				result := db.Exec(string(content))
				if result.Error != nil {
					return result.Error
				}
				log.Printf("Успешно выполнен SQL-скрипт: %s", info.Name())
			} else {
				// Хеш не изменился, пропускаем выполнение
				log.Printf("Скрипт %s не изменился, выполнение пропущено", info.Name())
			}
		}

		return nil
	})
}
