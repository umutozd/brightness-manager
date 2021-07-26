package brightness

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	ConfigFile            string   `json:"-"`
	File                  *os.File `json:"-"`
	Device                string   `json:"-"`
	LastAppliedDevice     string   `json:"last_applied_device,omitempty"`
	LastAppliedBrightness float64  `json:"last_applied_brightness,omitempty"`
	Increase              float64  `json:"-"`
	Decrease              float64  `json:"-"`
	Refresh               bool     `json:"-"`
}

func NewConfig() *Config {
	cfgFile := ""
	home, err := os.UserHomeDir()
	if err == nil {
		cfgFile = filepath.Join(home, ".config", "brightness", "config.json")
	}
	return &Config{
		LastAppliedDevice:     "HDMI-1-2",
		LastAppliedBrightness: 1.0,
		Device:                "HDMI-1-2",
		ConfigFile:            cfgFile,
		Increase:              0,
		Decrease:              0,
		Refresh:               false,
	}
}

func (cfg *Config) OpenConfigFile() error {
	file, err := os.OpenFile(cfg.ConfigFile, os.O_RDWR, 0766)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error openning config file: %v", err)
		}

		// file does not exist, open it with create flag
		file, err = os.OpenFile(cfg.ConfigFile, os.O_CREATE|os.O_RDWR, 0766)
		if err != nil {
			return fmt.Errorf("error openning config file with create flag: %v", err)
		}
		cfg.LastAppliedBrightness = 1.0
		cfg.LastAppliedDevice = cfg.Device
		cfg.File = file
	} else {
		// file exists and has been opened, read and unmarshal it
		cfg.File = file
		content, err := ioutil.ReadAll(cfg.File)
		if err != nil {
			return fmt.Errorf("error reading the content of the config file: %v", err)
		}
		if err = json.Unmarshal(content, cfg); err != nil {
			return fmt.Errorf("error unmarshaling the content of the config file: %v", err)
		}
	}
	return nil
}

func (cfg *Config) SaveConfigFile(brightness float64) error {
	var err error
	cfg.LastAppliedDevice = cfg.Device
	cfg.LastAppliedBrightness = brightness

	if _, err = cfg.File.Seek(0, 0); err != nil {
		return fmt.Errorf("error seeking config file's start: %v", err)
	}

	if err = cfg.File.Truncate(0); err != nil {
		return fmt.Errorf("error truncating config file while saving: %v", err)
	}

	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling while saving config file: %v", err)
	}
	if _, err = cfg.File.Write(content); err != nil {
		return fmt.Errorf("error writing to config file while saving: %v", err)
	}

	return nil
}
