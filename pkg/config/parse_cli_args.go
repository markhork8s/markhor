package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/markhork8s/markhor/pkg"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const versionHelpMessage = "Print the version of this program and exit"

func parseCliArgs() error {
	// Define CLI flags
	pflag.StringP("config", "c", pkg.DEFAULT_CONFIG_PATH, "Path to config file")
	helpSet := pflag.BoolP("help", "h", false, "Show this help message")
	versionSet := pflag.BoolP("version", "v", false, versionHelpMessage)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return errors.New("Could not parse CLI flags: " + err.Error())
	}

	// Print help or version message if needed
	if *helpSet {
		pflag.PrintDefaults()
		os.Exit(0)
	} else if *versionSet {
		fmt.Printf("v%s\n", pkg.VERSION)
		os.Exit(0)
	}

	return nil
}
