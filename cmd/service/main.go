package main

import (
	"fmt"
	"go-playground/internal/config"
	"go-playground/internal/logger"
	"go-playground/internal/repositories"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Hello World")

	// Initial Config
	cfg := config.InitConfig()

	// Initial Logger
	appLogger := logger.NewLogger(cfg.Log)
	defer func() {
		// appLogger.Logger.Sync() // ถ้าเป็น *zap.Logger โดยตรง
		if zl, ok := appLogger.(*logger.ZapLogger); ok {
			err := zl.Sync()
			if err != nil {
				return
			}
		}

	}()

	appLogger.Info("Logger initialized successfully.",
		zap.String("level", cfg.Log.Level),
		zap.String("format", cfg.Log.Format),
		zap.String("output", cfg.Log.Output))

	// Initial Database
	dbSql := repositories.NewMySQL(cfg.Database, appLogger)
	if dbSql == nil {
		appLogger.Fatal("Failed to initialize database connection.")
	}

	// ... โค้ดส่วนอื่นๆ ของแอปพลิเคชัน
	appLogger.Info("Application started successfully!")
}
