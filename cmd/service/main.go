package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go-playground/internal/api"
	"go-playground/internal/config"
	"go-playground/internal/logger"
	"go-playground/internal/repositories"
	"go.uber.org/zap"
	"net/http"
)

func main() {

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

	//Init Echo service API
	e := echo.New()
	router := api.NewRouter(e, cfg, appLogger)

	appLogger.Info("Starting Echo server...", zap.String("port", cfg.App.Port))
	if err := router.EchoInstance.Start(cfg.App.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		appLogger.Fatal("Echo server failed to start.", zap.Error(err))
	}

	appLogger.Info("Application gracefully shut down.")

}
