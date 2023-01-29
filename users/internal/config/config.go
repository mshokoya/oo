package config

import (
	"github.com/spf13/viper"
)


type Config struct {
	MONGO_URI string
	PORT int
	DB_NAME string
}

func LoadEnvVars[T any](path string, env *T) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	// viper.SetDefault("database.dbname", "test_db")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&env)
	if err != nil {
		return err
	}
	return nil
}