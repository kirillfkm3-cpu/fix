package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB открывает соединение с SQLite и применяет миграции
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Проверяем, что соединение реально установлено
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Настройка драйвера для миграций
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Ошибка драйвера миграций: %v", err)
	}

	// Указываем путь к папке с SQL файлами (относительно корня бэкенда)
	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/db/migrations/sqlite",
		"sqlite3", driver)
	if err != nil {
		log.Fatalf("Ошибка инициализации миграций: %v", err)
	}

	// Накатываем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	log.Println("База данных успешно подключена и миграции применены!")
	return db, nil
}