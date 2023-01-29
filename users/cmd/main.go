package main

import (
	"ecom-users/internal/application"
	"ecom-users/internal/config"
	"ecom-users/internal/logger"
	"ecom-users/internal/repository"
	"ecom-users/internal/server/http"
	"flag"
	"fmt"
	"os"
)

func main() {
	var env config.Config
	envPath := flag.String("env", "internal", "Path to env file")

	err := config.LoadEnvVars(*envPath, &env)
	if err != nil {
		fmt.Println("Env Error")
		panic(err.Error())
	}

	logger := logger.New(os.Stdout, logger.LevelInfo)
	db, err := repository.OpenDB(&env)

	dbModels := repository.NewModels(db)

	app := &application.Application{
		Logger: logger,
		Models: dbModels,
		Config: &env,
	}

	srv := http.New(app)
}