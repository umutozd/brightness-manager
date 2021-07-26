# Brightness Controller
This Go project is a shortcut for myself for easily `increasing`, `decreasing` and `refreshing` the brightness on any of my display devices.

## Installation
```bash
git clone git@github.com:umutozd/brightness-manager.git     # clone
cd brightness-manager                                       # enter downloaded directory
make                                                        # build
make copy-to-path                                           # sudo copy to /usr/bin/
```

## Usage
```bash
--config-file value         The absolute path to the config file. (default: "$HOME/.config/brightness/config.json")
-i value, --increase value  Increases the brightness by the given amount. (default: 0)
-d value, --decrease value  Decreases the brightness by the given amount. (default: 0)
--device value              The name of the device whose brightness is to be updated. (default: "HDMI-1-2")
-r, --refresh               When specified, the last applied brightness is applied.
--help, -h                  show help
--version, -v               print the version
```
## Tips
Add keyboard shortcuts from the settings of the desktop environment you're using for increasing and decreasing a monitor's brightness. Furthermore, to apply the last brightness configuration on your devices on boot-time, use brightness-manager with `--refresh` flag to re-apply the last configuration.

## Possible Improvements
- Adding support for dynamic device names.
    - This can be achieved by adding a bash auto-complete feature.
