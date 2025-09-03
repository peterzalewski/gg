package model

import (
	"errors"
	"os"
	"path"
)

type Repository struct {
	Worktree     string
	GitDirectory string
}

var ErrWorktreeMustBeDir = errors.New("worktree must be a directory")

func NewRepository(worktree string) (*Repository, error) {
	s, err := os.Stat(worktree)
	if err != nil {
		return nil, err
	}
	if !s.IsDir() {
		return nil, ErrWorktreeMustBeDir
	}

	return &Repository{
		GitDirectory: path.Join(worktree, ".git"),
		Worktree:     worktree,
	}, nil
}

func (r Repository) GitPath(elem ...string) string {
	foo := append([]string{r.GitDirectory}, elem...)
	return path.Join(foo...)
}
