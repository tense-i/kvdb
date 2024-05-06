package bitcask_go

import "bitcask-go/index"

type Options struct {
	DirPath         string          //数据库文件的目录
	DataFileMaxSize int64           //数据文件的最大大小
	SyncWrite       bool            //每次写操作后是否同步（持久化）
	indexType       index.IndexType //索引类型
}

type IndexType = int8

const (
	// Btree 索引
	bt IndexType = iota + 1
	//ART 自适应基数树索引
	ART
)
