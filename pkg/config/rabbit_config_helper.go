package config

import (
	config2 "github.com/exgamer/gosdk-core/pkg/config"
)

// InitRabbitConfig Инициализация конфига
func InitRabbitConfig() (*RabbitConfig, error) {
	dbConfig := &RabbitConfig{}
	err := config2.InitConfig(dbConfig)

	if err != nil {
		return nil, err
	}

	return dbConfig, nil
}
