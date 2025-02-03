/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/sks/kihocche/pkg/constants"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print the version number of %s", os.Args[0]),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "%s %s\n", os.Args[0], constants.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
