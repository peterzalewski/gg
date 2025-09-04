package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"petezalew.ski/pit/model"
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "Provide content or type and size information for repository objects",
	Run:   catFile,
}

func catFile(cmd *cobra.Command, args []string) {
	worktree := "."
	repo, err := model.NewRepository(worktree)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	o, err := repo.ReadObject(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(o)
}
