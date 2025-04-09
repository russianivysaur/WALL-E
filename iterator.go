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

var ErrCorruptedEntry error = errors.New("crc32 checksum does not match")

// Iterator
// Returns log entries as a struct
type Iterator struct {
	currentSegmentNumber uint64
	currentSegment       *os.File
	directory            string
}

var ErrEndOfLog error = errors.New("end of log")

func NewIterator(directory string) (*Iterator, error) {
	logFiles, err := filepath.Glob(filepath.Join(directory, fmt.Sprintf("%s*", segmentFilePrefix)))
	if err != nil {
		return nil, err
	}
	segmentFile, err := getFirstSegmentFile(logFiles, directory)
	return &Iterator{
		directory:            directory,
		currentSegmentNumber: 0,
		currentSegment:       segmentFile,
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
			return iterator.Next()
		}
	}
	size := binary.BigEndian.Uint64(sizeBytes)
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
	crcBuffer := make([]byte, 8)
	if _, err := binary.Encode(crcBuffer, binary.BigEndian, entry.Lsn); err != nil {
		return nil, err
	}
	if !verifyChecksum(append(crcBuffer, entry.Data...), &entry) {
		return nil, ErrCorruptedEntry
	}
	return &entry, nil
}

func (iterator *Iterator) jumpSegment() error {
	nextSegmentName := filepath.Join(iterator.directory, fmt.Sprintf("%s%d", segmentFilePrefix, iterator.currentSegmentNumber+1))
	file, err := os.OpenFile(nextSegmentName, os.O_RDONLY, 0777)
	if errors.Is(err, os.ErrNotExist) {
		return ErrEndOfLog
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	iterator.currentSegment = file
	iterator.currentSegmentNumber++
	return nil
}
