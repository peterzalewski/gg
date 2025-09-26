package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"petezalew.ski/gg/model"
)

func NewCatFileCmd() *cobra.Command {
	var catFileCmd = &cobra.Command{
		Use:   "cat-file",
		Short: "Provide content or type and size information for repository objects",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, ok := cmd.Context().Value(model.Repository{}).(*model.Repository)
			if !ok {
				return fmt.Errorf("could not retrieve repo from context")
			}

			o, err := repo.ReadObject(args[0])
			if err != nil {
				return err
			}

			fmt.Print(o)
			return nil
		},
	}
	return catFileCmd
}
