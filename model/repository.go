package model

import (
	"errors"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

type Repository struct {
	Worktree     string
	GitDirectory string
}

var ErrWorktreeMustBeDir = errors.New("worktree must be a directory")
var ErrNotAGitRepository = errors.New("not a git repository")
var ErrBadRevision = errors.New("bad revision")
var indirectRefRe = regexp.MustCompile(`^ref: (?P<indirectRef>[^\n]+)`)
var hashRe = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)

type RepositoryOption func(*Repository) error

func WithDiscoverRoot() RepositoryOption {
	return func(r *Repository) error {
		worktree, err := os.Getwd()
		if err != nil {
			return err
		}

		for {
			gitDirectory := path.Join(worktree, ".git")

			s, err := os.Stat(gitDirectory)
			if err == nil && s.IsDir() {
				r.Worktree = worktree
				r.GitDirectory = gitDirectory
				return nil
			}

			if worktree == "/" {
				return ErrNotAGitRepository
			}

			worktree = path.Dir(worktree)
		}
	}
}

func WithRoot(worktree string) RepositoryOption {
	return func(r *Repository) error {
		if wt, err := os.Stat(worktree); err != nil || !wt.IsDir() {
			return ErrWorktreeMustBeDir
		}

		gitDirectory := path.Join(worktree, ".git")
		if gd, err := os.Stat(gitDirectory); err != nil || !gd.IsDir() {
			return ErrNotAGitRepository
		}

		r.GitDirectory = gitDirectory
		r.Worktree = worktree

		return nil
	}
}

func NewRepository(options ...RepositoryOption) (*Repository, error) {
	r := &Repository{}

	for _, option := range options {
		err := option(r)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r Repository) GitPath(elem ...string) string {
	foo := append([]string{r.GitDirectory}, elem...)
	return path.Join(foo...)
}

func (r Repository) ResolveRef(filename string) (string, error) {
	// If it looks like a full SHA1, check it exists and, if it does,
	// we'll lazily assume it's a commit and return it
	if hashRe.MatchString(filename) {
		path := r.GitPath("objects", filename[:2], filename[2:])
		_, err := os.Stat(path)
		if err != nil {
			return "", err
		}

		return filename, nil
	}

	// Otherwise, try to open it and see if it's an indirect reference
	// `ref: refs/heads/blah` or a SHA
	path := r.GitPath(filename)
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	contents, _ := io.ReadAll(file)
	if match := indirectRefRe.FindSubmatch(contents); match != nil {
		filename = string(match[indirectRefRe.SubexpIndex("indirectRef")])
		return r.ResolveRef(filename)
	}

	return r.ResolveRef(strings.TrimSpace(string(contents)))
}
