package loconfig

import (
	"fmt"
	"ms-teacher/api/constants"

	"github.com/Muraddddddddd9/ms-database/config"
	"github.com/joho/godotenv"
)

type LocalConfig struct {
	ORIGIN_URL   string
	PROJECT_PORT string
	NGINX_URL    string

	ADMIN_EMAIL    string
	ADMIN_PASSWORD string
}

func LoadLocalConfig() (*LocalConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrLoadEnv)
	}

	return &LocalConfig{
		ORIGIN_URL:   config.GetEnv("ORIGIN_URL"),
		PROJECT_PORT: config.GetEnv("PROJECT_PORT"),
		NGINX_URL:    config.GetEnv("NGINX_URL"),
	}, nil
}
