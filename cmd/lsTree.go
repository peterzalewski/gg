package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"petezalew.ski/pit/model"
)

var hashRe = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Short: "List the contents of a tree object",
	RunE:  lsTree,
}

func lsTree(cmd *cobra.Command, args []string) error {
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

	if !hashRe.MatchString(ref) {
		ref = repo.ResolveRef(ref)
	}

	o, err := repo.ReadObject(ref)
	if err != nil {
		return err
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
		return fmt.Errorf("invalid object type")
	}

	fmt.Println(tree)
	return nil
}
