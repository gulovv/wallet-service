package config

import (
    "fmt"
    "log"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    DBConnString string
    ServerAddress string
}

func LoadConfig() Config {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dbConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_SSLMODE"), 
    )

    return Config{
        DBConnString: dbConnString,
        ServerAddress: ":" + os.Getenv("PORT"),
    }
}