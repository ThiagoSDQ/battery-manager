package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setThresholdCmd = &cobra.Command{
	Use:   "set-threshold [value]",
	Short: "Sets the upper threshold for battery charging",
	Long:  `Sets the threshold value for when the battery should stop charging.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("set-threshold accepts one arg: value [1-100]")
			return
		}

		value, err := strconv.Atoi(args[0])
		if err != nil || value < 1 || value > 100 {
			fmt.Printf("invalid value: %s, value should be from 1 to 100\n", args[0])
			return
		}

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.Set("charging-threshold", value)
		viper.WriteConfig()
		exec.Command("bash", "-c", "echo "+viper.GetString("charging-threshold")+" > "+home+"/.battery-manager/.charging-threshold").Run()
		exec.Command("bash", "-c", "sudo systemctl start battery-manager.service").Run()
	},
}

func init() {
	rootCmd.AddCommand(setThresholdCmd)
}
