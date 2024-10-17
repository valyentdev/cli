package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func StoreAPIKey(key string) error {
	_, err := initConfigDir()
	if err != nil {
		return err
	}

	viper.Set("key", key)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func initConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(homeDir, ".valyent")
	if !directoryExists(dir) {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return "", err
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	if _, err := os.Stat(filepath.Join(dir, "config.yaml")); os.IsNotExist(err) {
		if _, err := os.Create(filepath.Join(dir, "config.yaml")); err != nil {
			return "", err
		}
	}

	return dir, nil
}

func directoryExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func RemoveConfigFile() error {
	configDir, err := initConfigDir()
	if err != nil {
		return err
	}

	return os.Remove(filepath.Join(configDir, "config.yaml"))
}