// config/config.go
package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DBUser     string
    DBPassword string
    DBName     string
    DBHost     string
    DBPort     string
}

func GetConfig() Config {
    // 加载.env文件
    err := godotenv.Load()
    if err != nil {
        log.Println("没有找到 .env 文件，使用环境变量或默认配置")
    }

    return Config{
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),
        DBHost:     os.Getenv("DB_HOST"),
        DBPort:     os.Getenv("DB_PORT"),
    }
}