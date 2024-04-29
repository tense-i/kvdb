package index

import (
	"bitcask-go/data"
	"github.com/google/btree"
	"sync"
)

// Btree 封装Google的btree库
type Btree struct {
	tree *btree.BTree
	lock *sync.RWMutex
}

func NewBtree() *Btree {
	return &Btree{
		tree: btree.New(32),
		lock: new(sync.RWMutex),
	}
}

func (bt *Btree) Put(key []byte, pos *data.LogRecordPos) bool {
	item := &Item{key: key, pos: pos}

	//加锁、Google的Btree不是并发安全的
	bt.lock.Lock()
	bt.tree.ReplaceOrInsert(item)
	bt.lock.Unlock()
	return true
}

func (bt *Btree) Get(key []byte) *data.LogRecordPos {
	item := &Item{
		key: key,
	}
	btreeItem := bt.tree.Get(item)
	if btreeItem == nil {
		return nil
	}
	return btreeItem.(*Item).pos
}
func (bt *Btree) Delete(key []byte) bool {
	it := &Item{key: key}
	olditem := bt.tree.Delete(it)
	if olditem != nil {
		return true
	}
	return false
}
