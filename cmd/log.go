package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"petezalew.ski/pit/model"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Run:   log,
}

func log(cmd *cobra.Command, args []string) {
	repo, err := model.NewRepository(model.WithDiscoverRoot())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var ref string
	if len(args) == 0 {
		ref = "HEAD"
	} else {
		ref = args[0]
	}

	hash := repo.ResolveRef(ref)
	for {
		o, err := repo.ReadObject(hash)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		commit, isCommit := o.(*model.Commit)
		if !isCommit {
			fmt.Printf("read something that isn't a commit: %s\n", hash)
			os.Exit(1)
		}

		message := commit.Message()
		newlineIndex := strings.Index(message, "\n")
		var firstLine string
		if newlineIndex == -1 {
			firstLine = message
		} else {
			firstLine = message[:newlineIndex]
		}
		fmt.Printf("*  %s | %s [%s]\n", hash[:7], firstLine, commit.Authors())

		parents := commit.Parents()
		if len(parents) == 0 {
			break
		}
		hash = parents[0]
	}
}
