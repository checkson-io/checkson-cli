package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"os"
	"path/filepath"
	"strings"
)

var cfgFile string
var Verbose bool
var DevMode bool

var envMapping = map[string]string{
	"REQUESTTIMEOUT": "CONTEXTS_DEFAULT_REQUESTTIMEOUT",
}

var configPaths = []string{"$HOME/.config/checkson", "$HOME/.checkson", "$SNAP_REAL_HOME/.config/checkson", "$SNAP_DATA/checkson", "/etc/checkson"}

func NewChecksonCommand(streams output.IOStreams) *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:   "checkson",
		Short: "command-line interface for Checkson",
		Long:  `A command-line interface for managing Checkson checks`,
	}

	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(newCompletionCmd())
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newListChecksCmd())
	rootCmd.AddCommand(newStatusChecksCmd())
	rootCmd.AddCommand(newListRunsCmd())
	rootCmd.AddCommand(newLoginCmd())
	rootCmd.AddCommand(newLogoutCmd())
	rootCmd.AddCommand(newLogsCmd())
	rootCmd.AddCommand(newCreateCheckCmd())
	rootCmd.AddCommand(newDeleteCheckCmd())
	rootCmd.AddCommand(newShowCheckCmd())
	rootCmd.AddCommand(newChannelCmd())

	// use upper-case letters for shorthand params to avoid conflicts with local flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config-file", "C", "", fmt.Sprintf("config file. one of: %v", configPaths))
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "V", false, "verbose output")
	rootCmd.PersistentFlags().Bool("dev-mode", false, "enable dev mode, communicating to local services")

	output.IoStreams = streams
	rootCmd.SetOut(streams.Out)
	rootCmd.SetErr(streams.ErrOut)
	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.Reset()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else if os.Getenv("CHECKSON_CONFIG") != "" {
		viper.SetConfigFile(os.Getenv("CHECKSON_CONFIG"))
	} else {
		for _, path := range configPaths {
			viper.AddConfigPath(os.ExpandEnv(path))
		}
		viper.SetConfigName("config")
	}

	if Verbose {
		output.IoStreams.EnableDebug()
	}

	if Verbose && os.Getenv("SNAP_NAME") != "" {
		output.Debugf("Running snap version %s on %s", os.Getenv("SNAP_VERSION"), os.Getenv("SNAP_ARCH"))
	}

	mapEnvVariables()

	replacer := strings.NewReplacer("-", "_", ".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("contexts.default.brokers", []string{"localhost:9092"})
	viper.SetDefault("current-context", "default")

	viper.SetConfigType("yml")
	viper.AutomaticEnv() // read in environment variables that match

	if err := readConfig(); err != nil {
		output.Fail(err)
	}
}

func mapEnvVariables() {
	for short, long := range envMapping {
		if os.Getenv(short) != "" && os.Getenv(long) == "" {
			_ = os.Setenv(long, os.Getenv(short))
		}
	}
}

func readConfig() error {
	var err error
	if err = viper.ReadInConfig(); err == nil {
		output.Debugf("Using config file: %s", viper.ConfigFileUsed())
		return nil
	}

	_, isConfigFileNotFoundError := err.(viper.ConfigFileNotFoundError)
	_, isOsPathError := err.(*os.PathError)

	if !isConfigFileNotFoundError && !isOsPathError {
		return errors.Errorf("Error reading config file: %s (%v)", viper.ConfigFileUsed(), err)
	} else {
		err = generateDefaultConfig()
		if err != nil {
			return errors.Wrap(err, "Error generating default config: ")
		}
	}

	// We read generated config now
	if err = viper.ReadInConfig(); err == nil {
		output.Debugf("Using config file: %s", viper.ConfigFileUsed())
		return nil
	} else {
		return errors.Errorf("Error reading config file: %s (%v)", viper.ConfigFileUsed(), err)
	}
}

// generateDefaultConfig generates default config in case there is no config
func generateDefaultConfig() error {

	cfgFile := filepath.Join(os.ExpandEnv(configPaths[0]), "config.yml")

	if os.Getenv("CHECKSON_CONFIG") != "" {
		// use config file provided via env
		cfgFile = os.Getenv("CHECKSON_CONFIG")
	} else if os.Getenv("SNAP_REAL_HOME") != "" {
		// use different configFile when running in snap
		for _, configPath := range configPaths {
			if strings.Contains(configPath, "$SNAP_REAL_HOME") {
				cfgFile = filepath.Join(os.ExpandEnv(configPath), "config.yml")
				break
			}
		}
	}

	if err := os.MkdirAll(filepath.Dir(cfgFile), os.FileMode(0700)); err != nil {
		return err
	}

	if err := viper.WriteConfigAs(cfgFile); err != nil {
		return err
	}

	output.Debugf("generated default config at %s", cfgFile)
	return nil
}
