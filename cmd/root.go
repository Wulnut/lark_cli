/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates and returns the root command for the application.
func NewRootCmd(deps Deps) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "lark",
		Short: "Lark CLI - Feishu Project command line tool",
		Long: `Lark CLI is a tool for interacting with the Feishu Project OpenAPI.
It provides a command line interface to manage and interact with various Feishu Project resources.`,
	}

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	deps := Deps{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	rootCmd := NewRootCmd(deps)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global and persistent flags can be defined here if needed.
}
