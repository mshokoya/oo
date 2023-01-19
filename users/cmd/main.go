package main

import (
	"ecom-users/internal/config"
	"ecom-users/internal/jsonLog"
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

	logger := jsonLog.New(os.Stdout, jsonLog.LevelInfo)
}