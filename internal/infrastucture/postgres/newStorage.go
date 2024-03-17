package postgres

import (
	"database/sql"
	"fmt"
	"time"
	"vk_test/internal/app"

	_ "github.com/lib/pq"
)

func NewStorage(cfg app.Config) (*Storage, error) {
	time.Sleep(5 * time.Second) // для корректного подключения в докере

	database_url := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.PG_HOST,
		cfg.PG_PORT,
		cfg.PG_DBNAME,
		cfg.PG_USER,
		cfg.PG_PASSWORD,
		cfg.PG_SSLMODE,
	)
	db, err := sql.Open("postgres", database_url)
	if err != nil {
		return nil, err
	}

	// проверка, что подключились
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}
