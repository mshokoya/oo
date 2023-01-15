package main

import (
	"ecom-users/internal/config"
	"flag"
)

type EnvVars struct {
	MONGO_URI string
	PORT int
}


func main() {
	var env EnvVars
	envPath := flag.String("env", "internal", "Path to env file")

	flag.Parse()

	err := config.LoadEnvVars(*envPath, &env)
	if err != nil {
		panic(err.Error())
	}


}