/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"app/log"
	"app/src"

	"github.com/spf13/cobra"
)

var activateC = &cobra.Command{
	Use:   "activate [groupName1] [groupName2] ...",
	Short: "Activate PIM for groups",
	Long:  `Activate PIM for groups`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.InitializeLogger()
		duration, _ := cmd.Flags().GetInt("duration")
		reason, _ := cmd.Flags().GetString("reason")
		interactive, _ := cmd.Flags().GetBool("interactive")
		headless, _ := cmd.Flags().GetBool("browser-headless")
		debug, _ := cmd.Flags().GetBool("debug")

		if headless && interactive {
			logger.Fatal("Cannot use headless and interactive flags at the same time")
		}

		if !headless && !interactive {
			logger.Warn("No flags were set; defaulting to headless mode.")
			headless = true
		}

		opts := src.ActivationOptions{
			Headless:    headless,
			Interactive: interactive,
			Reason:      reason,
			Duration:    duration,
			GroupNames:  args, // filter criteria for activation
			Debug:       debug,
		}
		src.ActivatePim(opts)
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(activateC)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// activateCmd.PersistentFlags().String("foo", "", "A help for foo")
	activateC.Flags().StringP("reason", "r", "", "Reason for activation")
	activateC.Flags().IntP("duration", "d", 8, "Duration of activation in hours")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// activateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	activateC.Flags().BoolP("interactive", "i", false, "If true will let you use browser to enter password")
}
