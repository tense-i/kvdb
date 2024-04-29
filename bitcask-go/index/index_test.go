package index

import (
	"bitcask-go/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBtree_Put(t *testing.T) {
	bt := NewBtree()
	res := bt.Put(nil, nil)
	assert.True(t, res)

	res = bt.Put([]byte("a"), &data.LogRecordPos{
		Fid:    1,
		Offset: 2,
	})

	assert.True(t, res)

}

func TestBtree_Get(t *testing.T) {
	bt := NewBtree()
	res1 := bt.Put(nil, &data.LogRecordPos{
		Fid:    1,
		Offset: 1,
	})
	assert.True(t, res1)

	pos1 := bt.Get(nil)
	assert.Equal(t, uint32(1), pos1.Fid)
	assert.Equal(t, int64(1), pos1.Offset)
	res1 = bt.Put([]byte("a"), &data.LogRecordPos{
		Fid:    1,
		Offset: 2,
	})
	assert.True(t, res1)

	pos1 = bt.Get([]byte("a"))
	assert.Equal(t, uint32(1), pos1.Fid)
	assert.Equal(t, int64(2), pos1.Offset)
}
