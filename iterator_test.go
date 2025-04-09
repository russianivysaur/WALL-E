package walle

import (
	"fmt"
	assertPkg "github.com/stretchr/testify/assert"
	"testing"
)

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

func TestMultipleSegments(t *testing.T) {
	assert := assertPkg.New(t)
	directory := "wal"
	var segmentSize int64 = 500

	wal, err := NewWal(directory, segmentSize)
	assert.NoError(err)

	recordCount := 100
	logs := make([][]byte, recordCount)
	for i := 0; i < recordCount; i++ {
		logs[i] = []byte(fmt.Sprintf("test log data %d", i))
		assert.NoError(wal.WriteLog(logs[i]))
	}
	iterator, err := NewIterator(directory)
	assert.NoError(err)
	for i := 0; i < recordCount; i++ {
		entry, err := iterator.Next()
		assert.NoError(err)
		assert.Equal(logs[i], entry.Data)
	}

	assert.NoError(clean(directory))
}
