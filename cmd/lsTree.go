package cmd

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"petezalew.ski/pit/model"
)

var hashRe = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Args:  cobra.MaximumNArgs(1),
	Short: "List the contents of a tree object",
	RunE:  lsTree,
}

// TODO: convert to NewLSTreeCmd() pattern so I can define flags here and not centrally in root.go

func printTree(hash string, repo *model.Repository, recursive bool, build *strings.Builder, parent string) error {
	o, err := repo.ReadObject(hash)
	if err != nil {
		return err
	}

	tree, ok := o.(*model.Tree)
	if !ok {
		return fmt.Errorf("object was not a tree: %s", hash)
	}

	for _, entry := range tree.Entries {
		o, err := repo.ReadObject(entry.Blob)
		if err != nil {
			return err
		}
		if o.ObjectType() == "tree" && recursive {
			printTree(entry.Blob, repo, recursive, build, path.Join(parent, entry.FileName))
		} else {
			build.WriteString(fmt.Sprintf("%6s %s %s     %s\n", entry.Mode, o.ObjectType(), entry.Blob, path.Join(parent, entry.FileName)))
		}
	}
	return nil
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

	var hash string
	switch o := o.(type) {
	case *model.Commit:
		hash = o.Tree
	case *model.Tree:
		hash = o.Hash
	default:
		return fmt.Errorf("invalid object type")
	}

	var build strings.Builder
	recursive, _ := cmd.Flags().GetBool("recursive")
	printTree(hash, repo, recursive, &build, "")
	fmt.Print(build.String())
	return nil
}
