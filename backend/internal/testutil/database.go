// Package testutil 提供无需外部数据库的 SQL 边界回归测试工具。
package testutil

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ysmmc/backend/internal/config"
	"github.com/ysmmc/backend/internal/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Database(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		DisableAutomaticPing: true, SkipDefaultTransaction: true, Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	previous, previousConfig := database.DB, config.AppConfig
	database.DB = db
	config.AppConfig = &config.Config{
		JWTSecret: "test-only-jwt-secret-at-least-32-characters", JWTExpireHours: 24,
		JWTRefreshExpireDays: 7, UploadPath: t.TempDir(), MaxDiskUsage: 100,
	}
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
		_ = sqlDB.Close()
		database.DB, config.AppConfig = previous, previousConfig
	})
	return mock
}
