package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"petezalew.ski/pit/model"
)

func NewLogCmd() *cobra.Command {
	var logCmd = &cobra.Command{
		Use:   "log",
		Short: "Show commit logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, ok := cmd.Context().Value(model.Repository{}).(*model.Repository)
			if !ok {
				return fmt.Errorf("could not retrieve repo from context")
			}

			var ref string
			if len(args) == 0 {
				ref = "HEAD"
			} else {
				ref = args[0]
			}

			hash, err := repo.ResolveRef(ref)
			if err != nil {
				return err
			}

			for {
				o, err := repo.ReadObject(hash)
				if err != nil {
					return err
				}

				commit, isCommit := o.(*model.Commit)
				if !isCommit {
					return fmt.Errorf("read something that isn't a commit: %s\n", hash)
				}

				fmt.Printf("*  %s | %s [%s]\n", hash[:7], commit.FirstLine(), commit.Authors)

				if len(commit.Parents) == 0 {
					break
				}
				hash = commit.Parents[0]
			}

			return nil
		},
	}
	return logCmd
}
