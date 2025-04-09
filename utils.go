package walle

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getLastSegmentFile(files []string, directory string) (*os.File, error) {
	highestSegmentNumber := 0
	for _, path := range files {
		_, filename := filepath.Split(path)
		fmt.Println(filename)
		segmentNumber, err := strconv.Atoi(strings.TrimPrefix(filename, segmentFilePrefix))
		if err != nil {
			return nil, err
		}
		highestSegmentNumber = max(highestSegmentNumber, segmentNumber)
	}
	filename := fmt.Sprintf("%s%d", segmentFilePrefix, highestSegmentNumber)
	filePath := filepath.Join(directory, filename)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return nil, err
	}
	return file, err
}

func getFirstSegmentFile(files []string, directory string) (*os.File, error) {
	highestSegmentNumber := 0
	for _, path := range files {
		_, filename := filepath.Split(path)
		segmentNumber, err := strconv.Atoi(strings.TrimPrefix(filename, segmentFilePrefix))
		if err != nil {
			return nil, err
		}
		highestSegmentNumber = min(highestSegmentNumber, segmentNumber)
	}
	filename := fmt.Sprintf("%s%d", segmentFilePrefix, highestSegmentNumber)
	filePath := filepath.Join(directory, filename)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	return file, err
}

func verifyChecksum(data []byte, entry *Entry) bool {
	return crc32.ChecksumIEEE(data) == entry.Checksum
}
