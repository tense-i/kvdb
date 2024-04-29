package index

import (
	"bitcask-go/data"
	"bytes"
	"github.com/google/btree"
)

// Indexer 抽象索引接口、后续如果需要接入其它数据结构，实现接口即可
type Indexer interface {
	// Put 插入Btree
	Put(key []byte, pos *data.LogRecordPos) bool
	// Get 根据key拿到记录的位置
	Get(key []byte) *data.LogRecordPos
	Delete(key []byte) bool
}

// Item 定义存储项、实现Google的Btree的item的接口
type Item struct {
	key []byte
	pos *data.LogRecordPos
}

func (ai *Item) Less(bi btree.Item) bool {
	//bi,(*item)断言
	return bytes.Compare(ai.key, bi.(*Item).key) == -1
}
