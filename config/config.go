package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var AppEnv string

type DbConnectionString struct {
	dbUser     string
	dbPassword string
	dbName     string
	dbHost     string
	dbPort     string
	dbSslMode  string
}

func GetConnectionString(configuration string) (string, error) {

	conStrData, err := loadConfigurationString(configuration)
	if err != nil {
		return "", fmt.Errorf("loadConfigurationString failed: %W", err)
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		conStrData.dbUser, conStrData.dbPassword, conStrData.dbName, conStrData.dbHost, conStrData.dbPort, conStrData.dbSslMode), nil
}

func loadConfigurationString(configuration string) (DbConnectionString, error) {
	user := os.Getenv("TEST_DB_USER")
	password := os.Getenv("TEST_DB_PASSWORD")
	name := os.Getenv("TEST_DB_NAME")
	sslMode := os.Getenv("TEST_DB_SSLMODE")

	var host, port string

	if configuration == "production" {
		host = os.Getenv("TEST_DB_HOST")
		port = os.Getenv("TEST_DB_PORT")
	} else {
		isDocker := os.Getenv("IS_DOCKER")
		if isDocker == "true" {
			fmt.Println("Using docker set up")
			host = os.Getenv("TEST_DB_HOST")
			port = "5432"
		} else {
			fmt.Println("Not using docker set up")
			host = "localhost"
			port = os.Getenv("TEST_DB_PORT")
		}
	}

	if user == "" || password == "" || name == "" || host == "" || port == "" || sslMode == "" {
		return DbConnectionString{}, errors.New("one or more required DB environmental variables are not set")
	}

	return DbConnectionString{
		dbUser:     user,
		dbPassword: password,
		dbName:     name,
		dbHost:     host,
		dbPort:     port,
		dbSslMode:  sslMode,
	}, nil
}

func InitGinMode() (string, error) {

	err := godotenv.Load()
	if err != nil {
		log.Println("could not load .env file, checking if local .env file exist")
	}

	AppEnv = os.Getenv("APP_ENV")

	//TODO: Implement usage of constants, to avoid usage of 'magic values'.
	if AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("Production config")
		return "production", nil
	} else {
		gin.SetMode(gin.DebugMode)
		fmt.Println("Development config")
		err = godotenv.Load("../.env")

		if err != nil {
			return "", errors.New("error loading .env file")
		}
		return "development", nil
	}
}
