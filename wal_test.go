package walle

import (
	assertPkg "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewWal(t *testing.T) {
	assert := assertPkg.New(t)
	directory := "wal"
	_, err := NewWal(directory, 500)
	assert.NoError(err)

	assert.NoError(clean(directory))
}

func TestWriteLog(t *testing.T) {
	assert := assertPkg.New(t)
	directory := "wal"
	var segmentSize int64 = 500
	wal, err := NewWal(directory, segmentSize)
	assert.NoError(err)
	testLog := "test log data"
	assert.NoError(wal.WriteLog([]byte(testLog)))

	assert.NoError(clean(directory))
}

func TestIterator(t *testing.T) {
	assert := assertPkg.New(t)
	directory := "wal"
	var segmentSize int64 = 500
	wal, err := NewWal(directory, segmentSize)
	assert.NoError(err)
	testLog := []byte("test log data")
	assert.NoError(wal.WriteLog(testLog))
	iterator, err := NewIterator(directory)
	assert.NoError(err)
	entry, err := iterator.Next()
	assert.NoError(err)
	assert.Equal(entry.Data, testLog)

	assert.NoError(clean(directory))
}

func clean(directory string) error {
	return os.RemoveAll(directory)
}
