package walle

import (
	"encoding/binary"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"path/filepath"
)

// Iterator
// Returns log entries as a struct
type Iterator struct {
	currentOffset  uint64
	currentSegment *os.File
}

var ErrEndOfLog error = errors.New("end of log")

func NewIterator(directory string) (*Iterator, error) {
	logFiles, err := filepath.Glob(filepath.Join(directory, fmt.Sprintf("%s*", segmentFilePrefix)))
	if err != nil {
		return nil, err
	}
	segmentFile, err := getFirstSegmentFile(logFiles, directory)
	return &Iterator{
		currentOffset:  0,
		currentSegment: segmentFile,
	}, nil
}

func (iterator *Iterator) Next() (*Entry, error) {
	sizeBytes := make([]byte, 8)
	if err := binary.Read(iterator.currentSegment, binary.BigEndian, &sizeBytes); err != nil {
		if errors.Is(err, io.EOF) {
			// reached end of segment
			// jump Segment
			if err := iterator.jumpSegment(); err != nil {
				return nil, err
			}
		}
	}
	size := binary.BigEndian.Uint64(sizeBytes)
	fmt.Printf("entry size %d", size)
	entryBytes := make([]byte, size)
	if n, err := iterator.currentSegment.Read(entryBytes); err != nil || uint64(n) != size {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(fmt.Sprintf("expected to read %d bytes for log entry but got %d bytes\n", size, n))
	}
	var entry Entry
	if err := proto.Unmarshal(entryBytes, &entry); err != nil {
		return nil, err
	}
	iterator.currentOffset++
	return &entry, nil
}

func (iterator *Iterator) jumpSegment() error {
	return nil
}
