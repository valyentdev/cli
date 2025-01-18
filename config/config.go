package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
	ravelAPI "github.com/valyentdev/ravel/api"
)

func StoreAPIKey(namespace, key string) (err error) {
	_, err = initConfigDir()
	if err != nil {
		return err
	}

	viper.Set("namespace", namespace)
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

func RetrieveNamespace() (string, error) {
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

	key := viper.GetString("namespace")

	return key, nil
}

func IsAlreadyInitialized() bool {
	vi := viper.New()
	vi.SetConfigName("valyent")
	vi.AddConfigPath(".")
	vi.SetConfigType("json")

	err := vi.ReadInConfig()

	return err == nil
}

func InitializeConfigFile(fleetID string) error {
	if _, err := os.Create("valyent.json"); err != nil {
		return err
	}

	vi := viper.New()
	vi.SetConfigName("valyent")
	vi.AddConfigPath(".")
	vi.SetConfigType("json")

	vi.Set("fleet_id", fleetID)

	if err := vi.WriteConfig(); err != nil {
		return err
	}

	return nil
}

type MachineTemplate struct {
	Name        string
	VCPUConfigs []VCPUConfig
}

type VCPUConfig struct {
	VCPUs         int
	MemoryConfigs []int
}

func AskForMachineSpecs() (*ravelAPI.GuestConfig, error) {
	templates := map[string][]VCPUConfig{
		"eco": {
			{VCPUs: 1, MemoryConfigs: []int{256, 512, 1024, 2048}},
			{VCPUs: 2, MemoryConfigs: []int{512, 1024, 2048, 4096}},
			{VCPUs: 4, MemoryConfigs: []int{1024, 2048, 4096, 8192}},
			{VCPUs: 8, MemoryConfigs: []int{2048, 4096, 8192, 16384}},
		},
		"std": {
			{VCPUs: 1, MemoryConfigs: []int{1024, 2048, 4096}},
			{VCPUs: 2, MemoryConfigs: []int{2048, 4096, 8192}},
			{VCPUs: 4, MemoryConfigs: []int{4096, 8192, 16384}},
		},
	}

	cfg := &ravelAPI.GuestConfig{}

	err := huh.NewForm(huh.NewGroup(
		// CPU kind selection
		huh.NewSelect[string]().
			Key("CpuKind").
			Title("CPU Kind").
			Height(4).
			Options(
				huh.NewOption("Eco", "eco"),
				huh.NewOption("Standard", "std"),
			).
			Value(&cfg.CpuKind),

		// VCPUs selection
		huh.NewSelect[int]().
			Title("VCPUs").
			Key("Cpus").
			Height(5).
			OptionsFunc(func() []huh.Option[int] {
				if cfg.CpuKind == "" {
					return []huh.Option[int]{}
				}
				var options []huh.Option[int]
				for _, config := range templates[cfg.CpuKind] {
					options = append(options, huh.NewOption(fmt.Sprintf("%d VCPUs", config.VCPUs), config.VCPUs))
				}
				return options
			}, cfg.CpuKind).
			Value(&cfg.Cpus),

		// Memory selection
		huh.NewSelect[int]().
			Key("MemoryMB").
			Title("RAM").
			Height(5).
			OptionsFunc(func() []huh.Option[int] {
				if cfg.CpuKind == "" || cfg.Cpus == 0 {
					return []huh.Option[int]{}
				}
				var options []huh.Option[int]
				for _, config := range templates[cfg.CpuKind] {
					if config.VCPUs == cfg.Cpus {
						for _, mem := range config.MemoryConfigs {
							options = append(options, huh.NewOption(fmt.Sprintf("%d MB", mem), mem))
						}
					}
				}
				return options
			}, &cfg.Cpus).
			Value(&cfg.MemoryMB),
	)).Run()
	if err != nil {
		return nil, fmt.Errorf("error collecting machine specs: %w", err)
	}

	return cfg, nil
}

type ProjectConfiguration struct {
	ravelAPI.CreateMachinePayload

	FleetID string `json:"fleet_id"`
}

func RetrieveProjectConfiguration() (cfg ProjectConfiguration, err error) {
	f, err := os.ReadFile("valyent.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(f, &cfg)
	if err != nil {
		return
	}

	return
}
