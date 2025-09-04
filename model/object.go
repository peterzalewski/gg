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
	hash     string
	oType    string
	size     int
	contents string
}

func (o Object) String() string {
	return o.contents
}

var ErrInvalidObject = errors.New("invalid git object")

func (r Repository) ReadObject(hash string) (*Object, error) {
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

	return &Object{
		hash:     hash,
		size:     size,
		oType:    string(contents[0:typeEnd]),
		contents: string(contents[headerEnd+1:]),
	}, nil
}
