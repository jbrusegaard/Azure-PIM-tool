/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ezpim",
	Short: "A CLI tool to simplify and automate PIM in Azure environments",
	Long: `This tool provides a CLI that simplifies and automates PIM in Azure environments.

You can authenticate with your own Azure account, list eligible roles, and easily activate one or more roles with configurable duration and reasons.
Uses the Azure CLI along with Playwright and Chromium under the hood. Run with the --help flag for more information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.UsageString())
	},
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolP(
		"browser-headless", "b", false,
		"browser headless, only use if you want the browser to be totally headless. All information will be taken via cli",
	)
	rootCmd.PersistentFlags().BoolP("debug", "", false, "Enable debug logging and errors")
}
