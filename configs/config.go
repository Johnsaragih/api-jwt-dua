package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	DB     DBConfig
	JWT    JWTConfig
}

type DBConfig struct {
	User    string
	Pass    string
	Host    string
	Port    string
	Name    string
	Charset string
}

type JWTConfig struct {
	Secret string
}

//variabe global config

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Tidak Ada File .env")
		log.Println("JWT SECRET:", os.Getenv("JWT_SECRET"))
	}

	AppConfig = &Config{
		AppEnv: os.Getenv("APP_ENV"),
		DB: DBConfig{
			User:    os.Getenv("DB_USER"),
			Pass:    os.Getenv("DB_PASS"),
			Host:    os.Getenv("DB_HOST"),
			Port:    os.Getenv("DB_PORT"),
			Name:    os.Getenv("DB_NAME"),
			Charset: os.Getenv("DB_CHAR"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
	}
}
