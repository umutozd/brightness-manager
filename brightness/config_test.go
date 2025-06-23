package brightness

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenConfigFileCreatesDir(t *testing.T) {
	root := t.TempDir()
	cfgPath := filepath.Join(root, "newdir", "config.json")

	if _, err := os.Stat(filepath.Dir(cfgPath)); !os.IsNotExist(err) {
		t.Fatalf("expected directory to not exist before call: %v", err)
	}

	cfg := &Config{ConfigFile: cfgPath, Device: "HDMI-1-2"}
	if err := cfg.OpenConfigFile(); err != nil {
		t.Fatalf("OpenConfigFile returned error: %v", err)
	}
	defer cfg.File.Close()

	if _, err := os.Stat(filepath.Dir(cfgPath)); err != nil {
		t.Fatalf("expected directory to exist after call: %v", err)
	}
}
