package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"petezalew.ski/pit/model"
)

var rootCmd = &cobra.Command{
	Use:   "pit",
	Short: "Do what git does in Go",
	Long:  "This is a toy version of git and a way to get familiar and comfortable with Go.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// TODO: bail if init?
		var rootOption model.RepositoryOption
		if workTree, err := cmd.Flags().GetString("work-tree"); err == nil {
			rootOption = model.WithRoot(workTree)
		} else {
			rootOption = model.WithDiscoverRoot()
		}

		repo, err := model.NewRepository(rootOption)
		if err != nil {
			return err
		}

		cmd.SetContext(context.WithValue(cmd.Context(), model.Repository{}, repo))
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("work-tree", "C", ".", "Run as if pit was started in this path instead of the current working directory.")

	rootCmd.AddCommand(catFileCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(lsTreeCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
