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
	envPath := flag.String("env", "env", "Path to env file")

	err := config.LoadEnvVars(*envPath, &env)
	if err != nil {
		fmt.Println("Env Error")
		panic(err.Error())
	}

	logger := logger.New(os.Stdout, logger.LevelInfo)

	db, err := repository.OpenDB(&env)
	if err != nil {
		fmt.Println("DB Error")
		panic(err.Error())
	}

	app := &application.Application{
		Logger: logger,
		Models: repository.NewModels(db),
		Config: &env,
	}

	srv := http.New(app)
	srv.Routes(app.Routes())
	srv.Serve(app)
}