package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pit",
	Short: "Do what git does in Go",
	Long:  "This is a toy version of git and a way to get familiar and comfortable with Go.",
}

func init() {
	rootCmd.AddCommand(catFileCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(logCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
