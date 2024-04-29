package fio

const DataFilePerm = 0644

type IOManager interface {
	// Read 从文件中根据具体位置读取对应数据
	Read([]byte, int64) (int, error)
	// Write 写入数据到文件中
	Write([]byte) (int, error)

	// Sync 同步内存映射
	Sync() error

	// Close 关闭文件
	Close() error
}
