package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"petezalew.ski/pit/model"
)

var hashRe = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Short: "List the contents of a tree object",
	Run:   lsTree,
}

func lsTree(cmd *cobra.Command, args []string) {
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

	if !hashRe.MatchString(ref) {
		ref = repo.ResolveRef(ref)
	}

	o, err := repo.ReadObject(ref)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var tree *model.Tree
	switch o := o.(type) {
	case *model.Commit:
		// Looks kinda gross
		p, _ := repo.ReadObject(o.Tree)
		tree, _ = p.(*model.Tree)
	case *model.Tree:
		tree = o
	default:
		fmt.Println("invalid object type")
		os.Exit(1)
	}

	fmt.Println(tree)
}
