package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/umutozd/brightness-controller/brightness"
	cli "gopkg.in/urfave/cli.v1"
)

func Refresh() error {
	logrus.Info("Refreshing...")
	cmd := exec.Command(
		"xrandr", "--output", cfg.LastAppliedDevice,
		"--brightness", fmt.Sprint(cfg.LastAppliedBrightness),
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Refresh error: %v", err)
	}
	return nil
}

func Increase() error {
	newValue := cfg.LastAppliedBrightness + cfg.Increase
	logrus.Infof("Increasing brightness by %f, to %f", cfg.Increase, newValue)
	cmd := exec.Command(
		"xrandr", "--output", cfg.Device,
		"--brightness", fmt.Sprint(newValue),
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Increase error: %v", err)
	}

	// successful, save this config
	return cfg.SaveConfigFile(newValue)
}

func Decrease() error {
	newValue := cfg.LastAppliedBrightness - cfg.Decrease
	logrus.Infof("Decreasing brightness by %f, to %f", cfg.Decrease, newValue)
	cmd := exec.Command(
		"xrandr", "--output", cfg.Device,
		"--brightness", fmt.Sprint(newValue),
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Decrease error: %v", err)
	}

	// successful, save this config
	return cfg.SaveConfigFile(newValue)
}

func Run(ctx *cli.Context) error {
	var err error
	logrus.Info("Entered Run()")
	if err = cfg.OpenConfigFile(); err != nil {
		return err
	}

	// process the arguments
	if cfg.Refresh {
		return Refresh()
	}
	if cfg.Increase != 0 {
		return Increase()
	}
	if cfg.Decrease != 0 {
		return Decrease()
	}

	// no actionable flags were provided
	cli.ShowAppHelp(ctx)
	return nil
}

var cfg = brightness.NewConfig()

func main() {
	app := cli.NewApp()
	app.Authors = []cli.Author{{Name: "Umut Özdoğan", Email: "umut.ozdgan@gmail.com"}}
	app.Description =
		`brightness-controller is a simple command-line tool that helps adjusting
		the brightness of multiple displays`
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config-file",
			Destination: &cfg.ConfigFile,
			Value:       cfg.ConfigFile,
			Usage:       "The absolute path to the config file.",
		},
		cli.StringFlag{
			Name:        "device",
			Destination: &cfg.Device,
			Value:       cfg.Device,
			Usage:       "The name of the device whose brightness is to be updated.",
		},
		cli.BoolFlag{
			Name:        "r, refresh",
			Destination: &cfg.Refresh,
			Usage:       "When specified, the last applied brightness is applied.",
		},
		cli.Float64Flag{
			Name:        "i, increase",
			Destination: &cfg.Increase,
			Value:       cfg.Increase,
			Usage:       "Increases the brightness by the given amount. Ignored if refresh flag is specified.",
		},
		cli.Float64Flag{
			Name:        "d, decrease",
			Destination: &cfg.Decrease,
			Value:       cfg.Decrease,
			Usage:       "Decreases the brightness by the given amount. Ignored if refresh flag is specified or if increase flag is specified.",
		},
	}
	app.Action = Run
	app.After = func(ctx *cli.Context) error {
		return cfg.File.Close()
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("Fatal error")
	}
}
