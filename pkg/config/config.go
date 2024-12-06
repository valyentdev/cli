package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
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

func InitializeConfigFile(fleetID string) error {
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

type (
	MachineConfig struct {
		Image       string      `json:"image" toml:"image"`
		Guest       GuestConfig `json:"guest" toml:"guest"`
		Workload    Workload    `json:"workload" toml:"workload"`
		StopConfig  *StopConfig `json:"stop_config,omitempty" toml:"stop_config,omitempty"`
		AutoDestroy bool        `json:"auto_destroy,omitempty" toml:"auto_destroy,omitempty"`
	}

	GuestConfig struct {
		CpuKind  string `json:"cpu_kind" toml:"cpu_kind"`
		MemoryMB int    `json:"memory_mb" toml:"memory_mb" minimum:"1"`
		Cpus     int    `json:"cpus" minimum:"1" toml:"cpus"`
	}

	Workload struct {
		Restart RestartPolicyConfig `json:"restart,omitempty" toml:"restart"`
		Env     []string            `json:"env,omitempty" toml:"env"`
		Init    InitConfig          `json:"init,omitempty" toml:"init"`
	}

	InitConfig struct {
		Cmd        []string `json:"cmd,omitempty" toml:"cmd"`
		Entrypoint []string `json:"entrypoint,omitempty" toml:"entrypoint"`
		User       string   `json:"user,omitempty" toml:"user"`
	}

	RestartPolicy string

	RestartPolicyConfig struct {
		Policy     RestartPolicy `json:"policy,omitempty" toml:"policy"`
		MaxRetries int           `json:"max_retries,omitempty" toml:"max_retries"`
	}

	StopConfig struct {
		Timeout *int    `json:"timeout,omitempty" toml:"timeout"` // in seconds
		Signal  *string `json:"signal,omitempty" toml:"signal"`
	}
)

type ProjectConfiguration struct {
	Config    MachineConfig `json:"config" toml:"config"`
	FleetID   string        `json:"-" toml:"fleet_id"`
	SkipStart bool          `json:"skip_start" toml:"skip_start"`
	Region    string        `json:"region" toml:"region"`
}

func RetrieveProjectConfiguration() (cfg ProjectConfiguration, err error) {
	f, err := os.ReadFile("valyent.toml")
	if err != nil {
		return
	}

	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		return
	}

	return
}
