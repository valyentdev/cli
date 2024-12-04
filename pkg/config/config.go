package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func StoreAPIKey(key string) (err error) {
	_, err = initConfigDir()
	if err != nil {
		return err
	}

	viper.Set("key", key)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return
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

func RetrieveAPIKey() (string, error) {
	configDir, err := initConfigDir()
	if err != nil {
		return "", err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	key := viper.GetString("key")

	return key, nil
}

func IsAlreadyInitialized() bool {
	vi := viper.New()
	vi.SetConfigName("valyent")
	vi.AddConfigPath(".")
	vi.SetConfigType("toml")

	err := vi.ReadInConfig()

	return err == nil
}

func InitializeConfigFile(
	fleetID string,
) error {
	vi := viper.New()
	vi.SetConfigName("valyent")
	vi.AddConfigPath(".")
	vi.SetConfigType("toml")

	vi.Set("fleet_id", fleetID)

	var fileExists bool

	_, err := os.Stat("valyent.toml")
	if err == nil {
		fileExists = true
	} else if os.IsNotExist(err) {
		fileExists = false
	} else {
		return err
	}

	if fileExists {
		err = vi.MergeInConfig()
		if err != nil {
			return err
		}
	} else {
		_, err = os.Create("valyent.toml")
		if err != nil {
			return err
		}
	}

	err = vi.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}
