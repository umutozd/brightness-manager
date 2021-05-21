package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type config struct {
	File                  *os.File `json:"-"`
	LastAppliedDevice     string   `json:"last_applied_device,omitempty"`
	LastAppliedBrightness float64  `json:"last_applied_brightness,omitempty"`
	increase              bool
	decrease              bool
}

func initConfig() (*config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %v", err)
	}
	configFile := filepath.Join(home, ".config", "brightness", "config.json")

	file, err := os.OpenFile(configFile, os.O_RDWR, 0766)
	if err != nil {
		if os.IsNotExist(err) {
			// file does not exist, open it with CREATE flag
			file, err = os.OpenFile(configFile, os.O_CREATE|os.O_RDWR, 0766)
			if err != nil {
				return nil, fmt.Errorf("error opening config file with CREATE flag: %v", err)
			}
			return &config{
				File:                  file,
				LastAppliedDevice:     "HDMI-1-2",
				LastAppliedBrightness: 1,
			}, nil
		}
		// some other error
		return nil, fmt.Errorf("error opening config file: %v", err)
	}

	// file is open now, read its data and unmarshal into config
	data, err := ioutil.ReadAll(file)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	cfg := &config{
		File: file,
	}
	// the size of the file might be 0, use default values in this case
	if len(data) != 0 {
		if err = json.Unmarshal(data, cfg); err != nil {
			file.Close()
			return nil, fmt.Errorf("error unmarshaling config file: %v", err)
		}

	} else {
		cfg.LastAppliedBrightness = 1
		cfg.LastAppliedDevice = "HDMI-1-2"
	}

	return cfg, nil
}

func (cfg *config) runCommand() error {
	if cfg.increase {
		cfg.LastAppliedBrightness += 0.1
	} else if cfg.decrease {
		cfg.LastAppliedBrightness -= 0.1
	}
	cmd := exec.Command("xrandr", "--output", cfg.LastAppliedDevice, "--brightness", fmt.Sprint(cfg.LastAppliedBrightness))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("xrandr failed; error: %v, output: %s", err, string(output))
	}
	return nil
}

func (cfg *config) save() error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling: %v", err)
	}

	if err = cfg.File.Truncate(0); err != nil {
		return fmt.Errorf("error truncating config file to %d: %v", len(b), err)
	}

	if _, err = cfg.File.Seek(0, 0); err != nil {
		return fmt.Errorf("error seeking in config file: %v", err)
	}

	if _, err = cfg.File.Write(b); err != nil {
		return fmt.Errorf("error writing to config file: %v", err)
	}

	return nil

}

func (cfg *config) parseFlags() error {
	inc := flag.Bool("increase", false, "when specified, brightness is increased by 0.1")
	dec := flag.Bool("decrease", false, "when specified, brightness is decreased by 0.1")
	flag.Parse()

	if !*inc && !*dec {
		return fmt.Errorf("either --increase or --decrease flags must be specified")
	} else if *inc && *dec {
		return fmt.Errorf("both increase and decrease cannot be specified")
	}

	cfg.increase = *inc
	cfg.decrease = *dec
	return nil
}

func main() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("error initializing config: %v", err)
	}
	defer cfg.File.Close()

	if err = cfg.parseFlags(); err != nil {
		log.Fatalf("error parsing command line flags: %v", err)
	}

	if err = cfg.runCommand(); err != nil {
		log.Fatalf("error running command: %v", err)
	}

	if err = cfg.save(); err != nil {
		log.Fatalf("error saving config to file: %v", err)
	}

}
