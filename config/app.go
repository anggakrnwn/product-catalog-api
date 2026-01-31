package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	DBConn      string `mapstructure:"DB_CONN"`
}

func Load() Config {
	viper.AutomaticEnv()

	viper.BindEnv("PORT")
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("DB_CONN")

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	config := Config{
		Port:        port,
		Environment: viper.GetString("ENVIRONMENT"),
		DBConn:      viper.GetString("DB_CONN"),
	}

	if config.DBConn == "" {
		log.Println("Warning: DB_CONN is empty")
	}

	return config
}
