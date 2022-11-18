package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bambi2/balance-app/internal/config"
	"github.com/bambi2/balance-app/internal/delivery"
	"github.com/bambi2/balance-app/internal/repository"
	"github.com/bambi2/balance-app/internal/server"
	"github.com/bambi2/balance-app/internal/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Balance App API Miscroservice
// @version 1.0
// @description HTTP API Microservice for handling users' balance

// @host localhost:8000
// @BasePath /
// @schemes http
func main() {
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s\n", err.Error())
	}

	// if err := godotenv.Load(); err != nil {
	// 	logrus.Fatalf("error loading env variables: %s\n", err.Error())
	// }

	db, err := repository.NewPostgresDB(repository.Config{
		Host:         viper.GetString("db.host"),
		Port:         viper.GetString("db.port"),
		Username:     viper.GetString("db.username"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: viper.GetString("db.dbname"),
		SSLMode:      viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("error initilizing database: %s\n", err.Error())
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := delivery.NewHandler(service)
	server := server.NewServer()

	go func() {
		if err := server.Run(viper.GetString("port"), handler.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error occured while running http server: %s\n", err.Error())
		}
	}()

	logrus.Println("Balance App Started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logrus.Println("Balance App Shutting Down...")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured shutting down the server: %s\n", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured closing database connection: %s\n", err.Error())
	}
}
