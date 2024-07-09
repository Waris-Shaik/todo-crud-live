package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBPort     string
	DBName     string
}

var Envs Config

func init() {
	Envs = initConfig()
}

func initConfig() Config {

	godotenv.Load()

	port := os.Getenv("PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	nodeEnv := os.Getenv("NODE_ENV")

	if port == "" {
		log.Fatal("error: PORT environment variable is not set ")
	}
	if dbUser == "" {
		log.Fatal("error: DB_USER environment variable is not set")
	}
	if dbPassword == "" {
		log.Fatal("error: DB_PASSWORD environment variable is not set")
	}
	if dbAddress == "" {
		log.Fatal("error: DB_ADDRESS environment variable is not set")
	}
	if dbPort == "" {
		log.Fatal("error: DB_PORT environment variable is not set")
	}
	if dbName == "" {
		log.Fatal("error: DB_NAME environment variable is not set")
	}

	if jwtSecretKey == "" {
		log.Fatal("error: JWT_SECRET_KEY environment variable is not set")
	}

	if nodeEnv == "" {
		log.Fatal("error: NODE_ENV environment variable is not set")
	}

	dbAddress = fmt.Sprintf("%s:%s", dbAddress, dbPort)

	return Config{
		Port:       port,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBAddress:  dbAddress,
		DBPort:     dbPort,
		DBName:     dbName,
	}
}
