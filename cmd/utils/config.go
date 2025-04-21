package util

import "github.com/spf13/viper"

type Config struct {
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	KafkaBroker       string `mapstructure:"KAFKA_BROKER"`
	ProducerTopic     string `mapstructure:"PRODUCER_TOPIC"`
	ConsumerTopic     string `mapstructure:"CONSUMER_TOPIC"`
	ConsumerGroupID   string `mapstructure:"CONSUMER_GROUP_ID"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
