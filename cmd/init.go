package cmd

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	ini "gopkg.in/ini.v1"
	"petezalew.ski/gg/model"
)

var ErrGitExistsAndIsNotDir = errors.New(".git exists and is not a directory")
var ErrGitDirExistsAndIsNotEmpty = errors.New(".git exists and is not empty")

const (
	DefaultDescription = "Unnamed repository; edit this file 'description' to name the repository.\n"
	DefaultHead        = "ref: refs/heads/master\n"
)

func NewInitCmd() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Create a new git project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var worktree string
			if len(args) == 0 {
				worktree = "."
			} else {
				worktree = args[0]
			}

			repo, err := model.NewRepository(model.WithRoot(worktree))
			if err != nil {
				return err
			}

			syscall.Umask(0)

			s, err := os.Stat(repo.GitDirectory)
			if err != nil {
				// .git doesn't exist, so create it
				if err := os.Mkdir(repo.GitDirectory, 0755); err != nil {
					return err
				}
			} else if !s.IsDir() {
				// .git exists but it's a file; abort
				return ErrGitExistsAndIsNotDir
			} else {
				// .git is a directory...
				if ls, err := os.ReadDir(repo.GitDirectory); err != nil {
					return err
				} else if len(ls) > 0 {
					// ...but it's not empty; abort
					return ErrGitDirExistsAndIsNotEmpty
				}
			}

			for _, dir := range []string{"branches", "objects", "refs/tags", "refs/heads"} {
				if err := os.MkdirAll(repo.GitPath(dir), 0755); err != nil {
					return err
				}
			}

			description, err := os.OpenFile(repo.GitPath("description"), os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer description.Close()
			description.WriteString(DefaultDescription)

			head, err := os.OpenFile(repo.GitPath("HEAD"), os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer head.Close()
			head.WriteString(DefaultHead)

			config := ini.Empty()
			core := config.Section("core")
			core.Key("repositoryformatversion").SetValue("0")
			core.Key("filemode").SetValue("false")
			core.Key("bare").SetValue("false")

			c, err := os.OpenFile(repo.GitPath("config"), os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer c.Close()

			config.WriteToIndent(c, "        ")

			fmt.Println(repo)
			return nil
		},
	}
	return initCmd
}
