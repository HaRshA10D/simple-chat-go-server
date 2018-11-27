package utils

import (
	"errors"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"

	"github.com/spf13/viper"
)

func LoadConfig(fileName string, path string) (config *model.Config, outErr error) {
	v := viper.New()
	v.SetConfigType("json")
	v.SetConfigName(fileName)
	v.AddConfigPath(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.New("File not read")
	}
	if err := v.Unmarshal(&config); err != nil {
		return nil, errors.New("File not Unmarshal")
	}
	return config, nil
}
