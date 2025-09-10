package model

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	// "math"
)

// https://git-scm.com/docs/index-format
type IndexEntry struct {
	CTime_S     uint32 `offset:"0:4"`
	CTime_NS    uint32 `offset:"4:8"`
	MTime_S     uint32 `offset:"8:12"`
	MTime_NS    uint32 `offset:"12:16"`
	Dev         uint32 `offset:"16:20"`
	Inode       uint32 `offset:"20:24"`
	Mode        uint32 `offset:"24:28"`
	ObjectType  int    `offset:"24:26",mask:"7"`
	Permissions int    `offset:"24:26",mask:"511"`
	Uid         uint32 `offset:"26:30"`
	Gid         uint32 `offset:"30:34"`
	FileSize    uint32 `offset:"34:38"`
	Hash        string `offset:"40:60"`
	Flags       uint16 `offset:"60:62"`
	AssumeValid bool
	FileName    string
}

type Index struct {
	Version int
	Length  int
	Entries []IndexEntry
}

var ErrInvalidIndex = errors.New("invalid index format")

func NewIndex(r io.ReadSeeker) (*Index, error) {
	header := make([]byte, 12)
	n, err := r.Read(header)
	if err != nil || n != 12 || string(header[:4]) != "DIRC" {
		return nil, ErrInvalidIndex
	}

	index := &Index{
		Version: int(binary.BigEndian.Uint32(header[4:8])),
		Length:  int(binary.BigEndian.Uint32(header[8:])),
		Entries: make([]IndexEntry, 0),
	}

	for i := 0; i < index.Length; i++ {
		entryBytes := make([]byte, 62)
		_, err := r.Read(entryBytes)
		if err != nil {
			return nil, err
		}

		fileNameLength := binary.BigEndian.Uint16(entryBytes[60:62]) & 8191
		fileName := make([]byte, fileNameLength)
		_, err = r.Read(fileName)
		if err != nil {
			return nil, err
		}

		entry := &IndexEntry{
			Hash:     hex.EncodeToString(entryBytes[40:60]),
			FileName: string(fileName),
		}

		index.Entries = append(index.Entries, *entry)

		padding := int64(8 - ((fileNameLength + 62) % 8))
		r.Seek(padding, io.SeekCurrent)
	}

	return index, nil
}
