package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/umutozd/brightness-controller/brightness"
)

func setupTest(t *testing.T) (cleanup func()) {
	t.Helper()
	tmpDir := t.TempDir()

	// stub xrandr binary
	stub := filepath.Join(tmpDir, "xrandr")
	if err := ioutil.WriteFile(stub, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	cfg = brightness.NewConfig()
	cfg.ConfigFile = filepath.Join(tmpDir, "config.json")
	cfg.Device = "HDMI-1-2"
	if err := cfg.OpenConfigFile(); err != nil {
		t.Fatal(err)
	}

	return func() {
		cfg.File.Close()
		os.Setenv("PATH", oldPath)
	}
}

func captureLogs() (*bytes.Buffer, func()) {
	buf := new(bytes.Buffer)
	logger := logrus.StandardLogger()
	oldOut := logger.Out
	oldLevel := logger.Level
	logger.Out = buf
	logger.SetLevel(logrus.WarnLevel)
	return buf, func() {
		logger.Out = oldOut
		logger.SetLevel(oldLevel)
	}
}

func TestIncreaseClampsAboveOne(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()
	cfg.LastAppliedBrightness = 0.95
	cfg.Increase = 0.1

	buf, restore := captureLogs()
	defer restore()

	if err := Increase(); err != nil {
		t.Fatalf("Increase returned error: %v", err)
	}

	if cfg.LastAppliedBrightness != 1 {
		t.Errorf("expected brightness to be clamped to 1, got %v", cfg.LastAppliedBrightness)
	}

	if !strings.Contains(buf.String(), "clamping") {
		t.Errorf("expected warning about clamping, got logs: %s", buf.String())
	}
}

func TestDecreaseClampsBelowZero(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()
	cfg.LastAppliedBrightness = 0.05
	cfg.Decrease = 0.1

	buf, restore := captureLogs()
	defer restore()

	if err := Decrease(); err != nil {
		t.Fatalf("Decrease returned error: %v", err)
	}

	if cfg.LastAppliedBrightness != 0 {
		t.Errorf("expected brightness to be clamped to 0, got %v", cfg.LastAppliedBrightness)
	}

	if !strings.Contains(buf.String(), "clamping") {
		t.Errorf("expected warning about clamping, got logs: %s", buf.String())
	}
}
