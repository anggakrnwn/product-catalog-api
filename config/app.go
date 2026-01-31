package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	DBConn      string `mapstructure:"DB_CONN"`
}

func Load() Config {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	// set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "development")

	// config struct
	config := Config{
		Port:        viper.GetString("PORT"),
		Environment: viper.GetString("ENVIRONMENT"),
		DBConn:      viper.GetString("DB_CONN"),
	}

	// validasi
	if config.DBConn == "" {
		panic("DB_CONN is required in environment variables")
	}

	return config
}
