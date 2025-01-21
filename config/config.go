package config

import (
	"github.com/spf3/viper"
)

type Config struct {
	ServerAddress string
	MongoDB MongoDBConfig
	Redis RedisConfig
	JWTSecret string
	OTPExpiriation int
}

type MongoDBConfig struct {
	URI string
	Databasem string
}

type RedisConfig struct {
	Addr string
	Password string
	DB string
}

func Load ()(*Config, error){
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetConfigPath(".")

	if err := viper.ReadInConfig; err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Umarshal(&config); err != nil {
		return nil, err
	}

	return config&, nil
}