/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"app/src"

	"github.com/spf13/cobra"
)

var activateC = &cobra.Command{
	Use:   "activate <groupName>",
	Short: "Activate PIM for groups",
	Long:  `Activate PIM for groups`,
	Run: func(cmd *cobra.Command, args []string) {
		duration, _ := cmd.Flags().GetInt("duration")
		reason, _ := cmd.Flags().GetString("reason")
		opts := src.ActivationOptions{
			Reason:    reason,
			Duration:  duration,
			GroupName: args[0], // filter criteria for activation
		}
		src.ActivatePim(opts)
	},
	Args: cobra.ExactArgs(1),
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
}
