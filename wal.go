package walle

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// no log rotation here
// doesn't remove old logs

const segmentFilePrefix = "seg"

type Wal struct {
	lock                 sync.RWMutex
	directory            string
	segmentSize          int64
	currentSegment       *os.File
	currentSegmentWriter *bufio.Writer
	lastLSN              int64
	currentSegmentNumber int64
}

func NewWal(directory string, segmentSize int64) (*Wal, error) {
	if err := os.MkdirAll(directory, 0777); err != nil {
		return nil, err
	}
	segmentFiles, err := filepath.Glob(filepath.Join(directory, fmt.Sprintf("%s*", segmentFilePrefix)))
	var currentSegment *os.File
	var lastLSN int64
	var currentSegmentNumber int64
	if len(segmentFiles) == 0 {
		// make the first segment file
		if currentSegment, err = createNewSegmentFile(currentSegmentNumber, directory); err != nil {
			return nil, err
		}
	} else {
		// find last segment and last lsn
		lastSegmentFile, err := getLastSegmentFile(segmentFiles, directory)
		if err != nil {
			return nil, err
		}
		_, filename := filepath.Split(lastSegmentFile.Name())
		if currentSegmentNumber, err = strconv.ParseInt(strings.TrimPrefix(filename, segmentFilePrefix), 10, 64); err != nil {
			return nil, err
		}
		lastLSN, err = findLastLSNFromFile(lastSegmentFile)
		if err != nil {
			return nil, err
		}
		stat, err := lastSegmentFile.Stat()
		if err != nil {
			return nil, err
		}
		if stat.Size() == segmentSize {
			// create new segment file
			if currentSegment, err = createNewSegmentFile(currentSegmentNumber, directory); err != nil {
				return nil, err
			}
			currentSegmentNumber++
		} else {
			currentSegment = lastSegmentFile
		}
	}

	wal := &Wal{
		directory:            directory,
		segmentSize:          segmentSize,
		currentSegment:       currentSegment,
		currentSegmentWriter: bufio.NewWriter(currentSegment),
		lastLSN:              lastLSN,
		currentSegmentNumber: currentSegmentNumber,
	}
	return wal, nil
}

func createNewSegmentFile(number int64, directory string) (*os.File, error) {
	filename := fmt.Sprintf("%s%d", segmentFilePrefix, number)
	filePath := filepath.Join(directory, filename)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return nil, err
	}
	return file, nil
}

func findLastLSNFromFile(file *os.File) (int64, error) {
	var lastLSN int64
	sizeBytes := make([]byte, 8)
	for {
		if err := binary.Read(file, binary.BigEndian, &sizeBytes); err != nil {
			if errors.Is(err, io.EOF) {
				return lastLSN, nil
			} else {
				return 0, err
			}
		}
		size := binary.BigEndian.Uint64(sizeBytes)
		buffer := make([]byte, size)
		if _, err := file.Read(buffer); err != nil {
			return lastLSN, err
		}
		var entry Entry
		if err := proto.Unmarshal(buffer, &entry); err != nil {
			return lastLSN, err
		}
		lastLSN = int64(entry.GetLsn())
	}
}

func (wal *Wal) WriteLog(data []byte) error {
	wal.lock.Lock()
	defer wal.lock.Unlock()
	return wal.writeLog(data)
}

func (wal *Wal) writeLog(data []byte) error {
	lsn := wal.lastLSN + 1
	lsnBuffer := make([]byte, 8)
	_, err := binary.Encode(lsnBuffer, binary.BigEndian, lsn)
	if err != nil {
		return err
	}
	checksum := crc32.ChecksumIEEE(append(data, lsnBuffer...))
	entry := &Entry{
		Lsn:      lsn,
		Data:     data,
		Checksum: checksum,
	}
	entryBytes, err := proto.Marshal(entry)
	if err != nil {
		return err
	}
	size := uint64(len(entryBytes))
	var buffer *bytes.Buffer = bytes.NewBuffer(make([]byte, 0))
	sizeBytes := make([]byte, 8)
	if _, err := binary.Encode(sizeBytes, binary.BigEndian, size); err != nil {
		return err
	}
	buffer.Write(sizeBytes)
	buffer.Write(entryBytes)
	if err := wal.checkRotationAndRotate(entryBytes); err != nil {
		return err
	}
	return wal.write(buffer.Bytes())
}

func (wal *Wal) write(entryBytes []byte) error {
	n, err := wal.currentSegmentWriter.Write(entryBytes)
	if err != nil {
		return err
	}
	if err := wal.currentSegmentWriter.Flush(); err != nil {
		return err
	}
	if n != len(entryBytes) {
		return errors.New(fmt.Sprintf("expected to write %d bytes,wrote %d bytes", len(entryBytes), n))
	}
	return nil
}

func (wal *Wal) checkRotationAndRotate(entryBytes []byte) error {
	stat, err := wal.currentSegment.Stat()
	if err != nil {
		return err
	}
	if stat.Size()+int64(len(entryBytes)) > wal.segmentSize {
		return wal.rotate()
	}
	return nil
}

// just creates new log file
func (wal *Wal) rotate() error {
	file, err := createNewSegmentFile(wal.currentSegmentNumber+1, wal.directory)
	if err != nil {
		return err
	}
	if err = wal.currentSegmentWriter.Flush(); err != nil {
		return err
	}
	if err = wal.currentSegment.Sync(); err != nil {
		return err
	}
	wal.currentSegment = file
	wal.currentSegmentWriter = bufio.NewWriter(wal.currentSegment)
	wal.currentSegmentNumber++
	return nil
}
