package config

import "github.com/spf13/viper"

type Config struct {
	Port        string `mapstructure:"PORT"`
	UserSvcUrl  string `mapstructure:"USER_SVC_URL"`
	OrderSvcUrl string `mapstructure:"ORDER_SVC_URL"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./config/env")
	viper.SetConfigName("svc")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
