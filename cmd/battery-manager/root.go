package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	Port    int
	Name    string
	PathMap string `mapstructure:"path_map"`
}

var C config

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "battery-manager",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.battery-manager/.config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".battery-manager" (without extension).
		viper.AddConfigPath(home + "/.battery-manager")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in, else, creates default config.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		firstRunConfig()
	}
}

func firstRunConfig() {
	setDefaultConfigs()

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	exec.Command("bash", "-c", "mkdir "+home+"/.battery-manager").Run()

	viper.WriteConfigAs(home + "/.battery-manager/.config.yaml")
	viper.SetConfigFile(home + "/.battery-manager/.config.yaml")
	exec.Command("bash", "-c", "echo "+viper.GetString("charging-threshold")+" > "+home+"/.battery-manager/.charging-threshold").Run()

	serviceFile := fmt.Sprintf(`
[Unit]
Description=Battery Manager

[Service]
Type=oneshot
ExecStart=/bin/bash -c 'cat %s/.battery-manager/.charging-threshold > /sys/class/power_supply/BAT?/charge_control_end_threshold'

[Install]
WantedBy=multi-user.target suspend.target hibernate.target hybrid-sleep.target suspend-then-hibernate.target`, home)

	exec.Command("bash", "-c", "echo \""+serviceFile+"\" > "+home+"/.battery-manager/battery-manager.service").Run()
	exec.Command("bash", "-c", "sudo ln -s "+home+"/.battery-manager/battery-manager.service /etc/systemd/system/battery-manager.service").Run()
	exec.Command("bash", "-c", "sudo systemctl daemon-reload").Run()
	exec.Command("bash", "-c", "sudo systemctl enable battery-manager.service").Run()
	exec.Command("bash", "-c", "sudo systemctl start battery-manager.service").Run()
}

func setDefaultConfigs() {
	viper.Set("charging-threshold", 100)
}
