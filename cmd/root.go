package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"petezalew.ski/gg/model"
)

var rootCmd = &cobra.Command{
	Use:   "gg",
	Short: "Do what git does in Go",
	Long:  "This is a toy version of git and a way to get familiar and comfortable with Go.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		if cmd.Name() == "init" {
			return nil
		}

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
	rootCmd.PersistentFlags().StringP("work-tree", "C", ".", "Run as if gg was started in this path instead of the current working directory.")

	rootCmd.AddCommand(NewCatFileCmd())
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewLogCmd())
	rootCmd.AddCommand(NewLSTreeCmd())
	rootCmd.AddCommand(NewStatusCmd())
}

func Execute() error {
	return rootCmd.Execute()
}
