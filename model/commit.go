package model

import (
	"fmt"
	"regexp"
	"strings"
)

type Commit struct {
	*Object
}

var (
	parentRe  = regexp.MustCompile(`(?m)^parent (?P<hash>[a-f0-9]{40})$`)
	authorRe  = regexp.MustCompile(`(?m)^author (?P<author>.+)$`)
	treeRe    = regexp.MustCompile(`(?m)^tree (?P<hash>[a-f0-9]{40})$`)
	messageRe = regexp.MustCompile(`(?s)\n\n(?P<message>.+)$`)
)

func (c *Commit) Authors() []string {
	authors := make([]string, 0)
	for _, match := range authorRe.FindAllStringSubmatch(c.Contents, -1) {
		authors = append(authors, match[authorRe.SubexpIndex("author")])
	}
	return authors
}

func (c *Commit) Parents() []string {
	parents := make([]string, 0)
	for _, match := range parentRe.FindAllStringSubmatch(c.Contents, -1) {
		parents = append(parents, match[parentRe.SubexpIndex("hash")])
	}
	return parents
}

func (c *Commit) Tree() string {
	match := treeRe.FindStringSubmatch(c.Contents)
	return match[treeRe.SubexpIndex("hash")]
}

func (c *Commit) Message() string {
	match := messageRe.FindStringSubmatch(c.Contents)
	return match[messageRe.SubexpIndex("message")]
}

func (c *Commit) String() string {
	var build strings.Builder

	for _, author := range c.Authors() {
		build.WriteString(fmt.Sprintf("Author: %s\n", author))
	}

	for _, parent := range c.Parents() {
		build.WriteString(fmt.Sprintf("Parent: %s\n", parent))
	}

	build.WriteString(fmt.Sprintf("Tree: %s\n", c.Tree()))
	build.WriteString(fmt.Sprintf("Message: %s", c.Message()))

	return build.String()
}
