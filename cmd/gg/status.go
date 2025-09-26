package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"petezalew.ski/gg/internal/model"
)

func NewStatusCmd() *cobra.Command {
	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show the working tree status",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, ok := cmd.Context().Value(model.Repository{}).(*model.Repository)
			if !ok {
				return fmt.Errorf("could not retrieve repo from context")
			}

			indexReader, err := os.Open(repo.GitPath("index"))
			if err != nil {
				return err
			}
			defer indexReader.Close()

			status, err := model.NewIndex(indexReader)
			if err != nil {
				return err
			}

			var builder strings.Builder

			current, err := repo.CurrentBranch()
			if err != nil {
				return err
			}

			builder.WriteString(fmt.Sprintf("On branch %s\n", current))
			for _, entry := range status.Entries {
				builder.WriteString(fmt.Sprintf("%s %s\n", entry.Hash, entry.FileName))
			}
			fmt.Print(builder.String())
			return nil
		},
	}
	return statusCmd
}
