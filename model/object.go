package model

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"os"
	"strconv"
)

type Object struct {
	Hash     string
	Type     string
	Size     int
	Contents string
}

func (o Object) String() string {
	return o.Contents
}

var ErrInvalidObject = errors.New("invalid git object")

func (r Repository) ReadObject(hash string) (interface{}, error) {
	path := r.GitPath("objects", hash[:2], hash[2:])

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	z, _ := zlib.NewReader(file)
	contents, _ := io.ReadAll(z)

	headerEnd := bytes.Index(contents, []byte{0x00})
	if headerEnd == -1 {
		return nil, ErrInvalidObject
	}

	typeEnd := bytes.Index(contents[:headerEnd], []byte{' '})
	size, err := strconv.Atoi(string(contents[typeEnd+1 : headerEnd]))
	if err != nil {
		return nil, ErrInvalidObject
	}

	base := &Object{
		Hash:     hash,
		Size:     size,
		Type:     string(contents[0:typeEnd]),
		Contents: string(contents[headerEnd+1:]),
	}

	switch base.Type {
	case "blob":
		return NewBlob(base), nil
	case "commit":
		return NewCommit(base), nil
	case "tag":
		return NewTag(base), nil
	case "tree":
		return NewTree(base), nil
	default:
		return nil, ErrInvalidObject
	}
}
