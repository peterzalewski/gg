package model

import (
	"fmt"
	"regexp"
	"strings"
)

type Commit struct {
	*Object
	Authors []string
	Parents []string
	Tree    string
	Message string
}

var (
	parentRe  = regexp.MustCompile(`(?m)^parent (?P<hash>[a-f0-9]{40})$`)
	authorRe  = regexp.MustCompile(`(?m)^author (?P<author>.+)$`)
	treeRe    = regexp.MustCompile(`(?m)^tree (?P<hash>[a-f0-9]{40})$`)
	messageRe = regexp.MustCompile(`(?s)\n\n(?P<message>.+)$`)
)

func NewCommit(obj *Object) *Commit {
	c := &Commit{Object: obj}
	c.Authors = make([]string, 0)

	for _, match := range authorRe.FindAllSubmatch(c.Contents, -1) {
		c.Authors = append(c.Authors, string(match[authorRe.SubexpIndex("author")]))
	}

	c.Parents = make([]string, 0)
	for _, match := range parentRe.FindAllSubmatch(c.Contents, -1) {
		c.Parents = append(c.Parents, string(match[parentRe.SubexpIndex("hash")]))
	}

	match := treeRe.FindSubmatch(c.Contents)
	c.Tree = string(match[treeRe.SubexpIndex("hash")])

	match = messageRe.FindSubmatch(c.Contents)
	c.Message = string(match[messageRe.SubexpIndex("message")])

	return c
}

func (c *Commit) ObjectType() string {
	return "commit"
}

func (c *Commit) FirstLine() string {
	newlineIndex := strings.Index(c.Message, "\n")
	var firstLine string
	if newlineIndex == -1 {
		firstLine = c.Message
	} else {
		firstLine = c.Message[:newlineIndex]
	}
	return firstLine
}

func (c *Commit) String() string {
	var build strings.Builder

	for _, author := range c.Authors {
		build.WriteString(fmt.Sprintf("Author: %s\n", author))
	}

	for _, parent := range c.Parents {
		build.WriteString(fmt.Sprintf("Parent: %s\n", parent))
	}

	build.WriteString(fmt.Sprintf("Tree: %s\n", c.Tree))
	build.WriteString(fmt.Sprintf("Message: %s", c.Message))

	return build.String()
}
