package model

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

type Tree struct {
	*Object
	Entries []TreeEntry
}

type TreeEntry struct {
	Mode     string
	FileName string
	Blob     string
}

var (
	treeEntryRe = regexp.MustCompile(`(?m)(?P<mode>[0-9]{6}) (?P<filename>[^\x00]+)\x00(?P<hash>.{20})`)
)

// TODO: I'm not convinced regex is the way to go here but it's always the first tool I reach for and here we are
func NewTree(obj *Object) *Tree {
	tree := &Tree{Object: obj}
	tree.Entries = make([]TreeEntry, 0)
	for _, match := range treeEntryRe.FindAllSubmatch([]byte(obj.Contents), -1) {
		tree.Entries = append(tree.Entries, TreeEntry{
			Mode:     string(match[treeEntryRe.SubexpIndex("mode")]),
			FileName: string(match[treeEntryRe.SubexpIndex("filename")]),
			// TODO: Figure out why this matches 21 bytes instead of 20 on all but the first match
			Blob: hex.EncodeToString(match[treeEntryRe.SubexpIndex("hash")][:20]),
		})
	}
	return tree
}

func (t *Tree) String() string {
	var build strings.Builder

	for _, entry := range t.Entries {
		build.WriteString(fmt.Sprintf("%s %-40s%s\n", entry.Mode, entry.FileName, entry.Blob))
	}

	return build.String()
}
