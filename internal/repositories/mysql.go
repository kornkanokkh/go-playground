package repositories

import (
	"fmt"
	"go-playground/internal/config"
	"go-playground/internal/logger"
	"go.uber.org/zap"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormDatabase เป็น struct ที่ใช้สำหรับการจัดการ Error ของ GORM
type GormDatabase struct {
	DB *gorm.DB
}

func (db *GormDatabase) Error() error {
	return db.DB.Error
}

// NewMySQL ทำการเชื่อมต่อและคืนค่าอินสแตนซ์ของ GORM DB สำหรับ MySQL
func NewMySQL(dbCfg config.DatabaseConfig, log logger.Logger) *gorm.DB {
	var database *gorm.DB
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.DBName)

	for i := 1; i <= 3; i++ {
		database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Info("Database connection successful!",
				zap.String("db_host", dbCfg.Host),
				zap.String("db_name", dbCfg.DBName)) // log success with context
			break
		} else {
			log.Warn("Failed to connect to MySQL database. Retrying in 3 seconds...",
				zap.Int("attempt", i),
				zap.Error(err), // log an error object directly
				zap.String("db_host", dbCfg.Host),
				zap.String("db_name", dbCfg.DBName))
			time.Sleep(3 * time.Second)
		}
	}

	if err != nil {
		log.Fatal("Failed to connect to MySQL database after multiple retries.",
			zap.Error(err),
			zap.String("db_host", dbCfg.Host),
			zap.String("db_name", dbCfg.DBName))
	}

	return database
}
