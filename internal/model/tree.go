package model

import (
	"bytes"
	"encoding/hex"
	"fmt"
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

func NewTree(obj *Object) *Tree {
	tree := &Tree{Object: obj}
	tree.Entries = make([]TreeEntry, 0)
	data := []byte(tree.Contents)

	var spaceIndex, nulIndex, hashIndex int
	for len(data) > 0 {
		spaceIndex = bytes.IndexByte(data, ' ')
		if spaceIndex == -1 {
			break
		}
		mode := string(data[:spaceIndex])

		nulIndex = bytes.IndexByte(data, 0)
		if nulIndex == -1 {
			break
		}
		fileName := string(data[spaceIndex+1 : nulIndex])
		hashIndex = nulIndex + 1

		tree.Entries = append(tree.Entries, TreeEntry{
			Mode:     mode,
			FileName: fileName,
			Blob:     hex.EncodeToString(data[hashIndex : hashIndex+20]),
		})

		data = data[hashIndex+20:]
	}

	return tree
}

func (t *Tree) ObjectType() string {
	return "tree"
}

func (t *Tree) String() string {
	var build strings.Builder

	for _, entry := range t.Entries {
		build.WriteString(fmt.Sprintf("%s %-40s%s\n", entry.Mode, entry.FileName, entry.Blob))
	}

	return build.String()
}
