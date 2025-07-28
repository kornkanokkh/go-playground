// internal/logger/logger.go
package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"go-playground/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger Interface: กำหนดพฤติกรรมที่ Logger ควรมี
// นี่คือสิ่งที่ส่วนอื่นๆ ของแอปพลิเคชันจะเรียกใช้
type Logger interface {
	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Warn(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field) // Fatal จะเรียก os.Exit(1)
	Panic(message string, fields ...zap.Field) // Panic จะ panic
	With(fields ...zap.Field) Logger
}

// ZapLogger เป็น struct ที่ implement Logger interface
// และใช้ *zap.Logger เป็น core
type ZapLogger struct {
	*zap.Logger
}

// NewLogger สร้างและคืนค่า Logger instance ตามการตั้งค่า
func NewLogger(cfg config.LogConfig) Logger {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.Level))); err != nil {
		level = zapcore.InfoLevel // Default to Info if invalid
		fmt.Printf("Invalid log level '%s', defaulting to info.\n", cfg.Level)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // ISO 8601 for time
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // INFO, ERROR, etc.

	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else { // default to console
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var output io.Writer
	if strings.HasPrefix(cfg.Output, "file:") {
		filePath := strings.TrimPrefix(cfg.Output, "file:")
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to open log file %s, falling back to stdout: %v\n", filePath, err)
			output = os.Stdout
		} else {
			output = file
		}
	} else if cfg.Output == "stderr" {
		output = os.Stderr
	} else { // default to stdout
		output = os.Stdout
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(output), level)
	return &ZapLogger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))} // เพิ่ม stacktrace เมื่อเป็น Error ขึ้นไป
}

func (l *ZapLogger) Debug(message string, fields ...zap.Field) {
	l.Logger.Debug(message, fields...)
}

func (l *ZapLogger) Info(message string, fields ...zap.Field) {
	l.Logger.Info(message, fields...)
}

func (l *ZapLogger) Warn(message string, fields ...zap.Field) {
	l.Logger.Warn(message, fields...)
}

func (l *ZapLogger) Error(message string, fields ...zap.Field) {
	l.Logger.Error(message, fields...)
}

func (l *ZapLogger) Fatal(message string, fields ...zap.Field) {
	l.Logger.Fatal(message, fields...)
}

func (l *ZapLogger) Panic(message string, fields ...zap.Field) {
	l.Logger.Panic(message, fields...)
}

func (l *ZapLogger) With(fields ...zap.Field) Logger {
	return &ZapLogger{l.Logger.With(fields...)}
}
