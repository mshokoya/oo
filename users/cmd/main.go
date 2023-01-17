package main

import (
	"ecom-users/internal/config"
	"flag"
	"fmt"
)

type EnvVars struct {
	MONGO_URI string
	PORT int
}


func main() {
	var env EnvVars
	envPath := flag.String("env", "internal", "Path to env file")

	flag.Parse()

	fmt.Println("test dsddsdsadd3")

	err := config.LoadEnvVars(*envPath, &env)
	if err != nil {
		fmt.Println("Env Error")
		panic(err.Error())
	}
}