package config

import "github.com/spf13/viper"

type conf struct {
	MaxRequestsWithoutToken int    `mapstructure:"MAX_REQUESTS_WITHOUT_TOKEN_PER_SECOND"`
	MaxRequestsWithToken    int    `mapstructure:"MAX_REQUESTS_WITH_TOKEN_PER_SECOND"`
	TimeBlockInSecond       int    `mapstructure:"TIME_BLOCK_IN_SECOND"`
	RedisHost               string `mapstructure:"REDIS_HOST"`
	RedisPort               string `mapstructure:"REDIS_PORT"`
}

func LoadConfig(path string) (*conf, error) {

	var config *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		panic(err)
	}

	return config, err
}
