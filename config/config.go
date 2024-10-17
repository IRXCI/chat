package config

import (
	"github.com/joho/godotenv"
)

// Load ... Загружает конфигурацию из файла
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
