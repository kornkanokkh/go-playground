package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-playground/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// ----- API Handlers -----
// แยก Handler functions ออกไปในไฟล์ internal/api/handlers.go ก็ได้ เพื่อความเรียบร้อย

func (r *Router) healthCheckHandler(appLogger logger.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Ping database to confirm connection is active
		//sqlDB, err := dbRepo.DB.DB() // Get underlying *sql.DB from GORM
		//if err != nil {
		//	appLogger.Error("Health check failed: Could not get underlying DB.", zap.Error(err))
		//	return c.String(http.StatusInternalServerError, "Database connection error (ping failed)!")
		//}
		//if err = sqlDB.Ping(); err != nil {
		//	appLogger.Error("Health check failed: Database ping failed.", zap.Error(err))
		//	return c.String(http.StatusInternalServerError, "Database connection error (ping failed)!")
		//}

		appLogger.Info("Health check endpoint hit.", zap.String("path", "/health"))
		return c.String(http.StatusOK, "OK! Service is running and connected to DB.")
	}
}

func (r *Router) helloHandler(appLogger logger.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		message := fmt.Sprintf("Hello, %s!", name)
		appLogger.Info("Hello endpoint hit.", zap.String("name", name))
		return c.String(http.StatusOK, message)
	}
}

func (r *Router) getItemsHandler(appLogger logger.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		appLogger.Info("Get items endpoint hit.", zap.String("path", "/api/v1/items"))
		// ตัวอย่าง: ดึงข้อมูลจากฐานข้อมูล
		// var items []models.Item
		// if err := dbRepo.DB.Find(&items).Error; err != nil {
		// 	appLogger.Error("Failed to fetch items from DB.", zap.Error(err))
		// 	return c.String(http.StatusInternalServerError, "Failed to fetch items")
		// }
		// return c.JSON(http.StatusOK, items)
		return c.String(http.StatusOK, "List of items (not yet implemented)")
	}
}

func (r *Router) createItemHandler(appLogger logger.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		appLogger.Info("Get items endpoint hit.", zap.String("path", "/api/v1/items"))
		return c.String(http.StatusOK, "List of items (not yet implemented)")
	}
}
