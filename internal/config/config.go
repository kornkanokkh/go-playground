package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// Config เป็น struct หลักสำหรับเก็บค่า configuration ทั้งหมด
type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Log      LogConfig      `yaml:"log"`
}

// AppConfig เป็น struct สำหรับการตั้งค่าแอปพลิเคชัน
type AppConfig struct {
	Name string `yaml:"name"`
	Env  string `yaml:"env"`
	Port string `yaml:"port"`
}

// DatabaseConfig เป็น struct สำหรับการตั้งค่าฐานข้อมูล
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// LogConfig สำหรับการตั้งค่า Logger
type LogConfig struct {
	Level  string `yaml:"level"`  // debug, info, warn, error, fatal, panic
	Format string `yaml:"format"` // json, console
	Output string `yaml:"output"` // stdout, stderr, file:/path/to/file.log
}

// InitConfig เป็นฟังก์ชัน helper สำหรับเริ่มต้นการโหลด config
func InitConfig() *Config {
	configPath := GetConfigPath()
	cfg, err := LoadConfig(configPath)
	if err != nil {
		// Log fatal error ถ้าโหลด config ไม่ได้ เพราะเป็นส่วนสำคัญของแอปพลิเคชัน
		// หรือสามารถจัดการ error ได้ตามต้องการ
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}
	return cfg
}

// LoadConfig อ่านไฟล์ configuration จาก path ที่กำหนด
func LoadConfig(configPath string) (*Config, error) {
	// ตรวจสอบว่าไฟล์ config อยู่จริงหรือไม่
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s", configPath)
	}

	// อ่านเนื้อหาไฟล์ YAML
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	// Unmarshal ค่าจาก YAML ลงใน struct Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &cfg, nil
}

// GetConfigPath ช่วยในการกำหนด path ของ config file ได้อย่างยืดหยุ่น
// โดยจะพยายามหาไฟล์ config.yaml จากหลายตำแหน่ง
func GetConfigPath() string {
	// 1. ลองหาในโฟลเดอร์ config/ (เมื่อรันจาก root project)
	if _, err := os.Stat("config/config.yaml"); err == nil {
		return "config/config.yaml"
	}

	// 2. ลองหาใน Current Working Directory (สำหรับ Docker หรือการ deploy)
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml"
	}

	// 3. ลองหาเมื่อรันจาก subdirectory (เช่น ./internal/config)
	if _, err := os.Stat("../config/config.yaml"); err == nil {
		return "../config/config.yaml"
	}

	// 4. Fallback ไปยังตำแหน่งที่คาดการณ์ได้ แต่จะคืนค่าเป็น path และ LoadConfig จะแจ้ง error ถ้าไม่พบ
	return "config/config.yaml"
}
