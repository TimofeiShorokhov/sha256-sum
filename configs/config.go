package configs

import (
	"github.com/spf13/viper"
	"sha256-sum/repository"
)

func ParseConfig(dir string) (*repository.PostgresDB, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	var cfg *repository.PostgresDB
	err = viper.UnmarshalKey("db", &cfg)

	return cfg, nil
}
