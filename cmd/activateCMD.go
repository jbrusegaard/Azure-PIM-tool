/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"app/src"
	"fmt"

	"github.com/spf13/cobra"
)

var activateC = &cobra.Command{
	Use:   "activate <type> <filter>",
	Short: "Activate PIM for groups, resources, or roles",
	Long:  `Activate PIM for groups, resources, or roles`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("activate called")
		duration, _ := cmd.Flags().GetInt("duration")
		reason, _ := cmd.Flags().GetString("reason")
		opts := src.ActivationOptions{
			Reason:         reason,
			Duration:       duration,
			ActivationType: args[0], // type of activation, e.g., "group", "resource", "role"
			Filter:         args[1], // filter criteria for activation
		}
		src.ActivatePim(opts)
	},
	Args: cobra.ExactArgs(2),
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
