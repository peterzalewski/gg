package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"petezalew.ski/pit/model"
)

const (
	DefaultDescription = "Unnamed repository; edit this file 'description' to name the repository.\n"
	DefaultHead        = "ref: refs/heads/master\n"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new git project",
	Run:   initRepository,
}

func initRepository(cmd *cobra.Command, args []string) {
	var worktree string
	if len(args) == 0 {
		worktree = "."
	} else {
		worktree = args[0]
	}

	repo, err := model.NewRepository(worktree)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	syscall.Umask(0)

	s, err := os.Stat(repo.GitDirectory)
	if err != nil {
		// .git doesn't exist, so create it
		if err := os.Mkdir(repo.GitDirectory, 0755); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if !s.IsDir() {
		// .git exists but it's a file; abort
		fmt.Println(".git exists and is not a directory")
		os.Exit(1)
	} else {
		// .git is a directory...
		if ls, err := os.ReadDir(repo.GitDirectory); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else if len(ls) > 0 {
			// ...but it's not empty; abort
			fmt.Println(".git exists and is not empty")
			os.Exit(1)
		}
	}

	for _, dir := range []string{"branches", "objects", "refs/tags", "refs/heads"} {
		if err := os.MkdirAll(repo.GitPath(dir), 0755); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	description, err := os.OpenFile(repo.GitPath("description"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer description.Close()
	description.WriteString(DefaultDescription)

	head, err := os.OpenFile(repo.GitPath("HEAD"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer head.Close()
	head.WriteString(DefaultHead)

	fmt.Println(repo)
}
