package db

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/demurk/tinyurl/internal/config"
	"github.com/demurk/tinyurl/internal/logger"
)

var (
	conn     *sql.DB
	syncOnce sync.Once
)

func Connect() error {
	zapLogger := logger.GetLogger()
	var err error
	syncOnce.Do(func() {
		conn, err = sql.Open("postgres", *config.DBConnectionURL)
		if err != nil {
			return
		}
		err = conn.Ping()
		if err != nil {
			return
		}
		zapLogger.Info("Successfully connected to DB!")
	})
	if err != nil {
		zapLogger.Error("Failed to connect to DB", zap.Error(err))
		conn.Close()
	}
	return err
}

func GetConnection() *sql.DB {
	return conn
}
