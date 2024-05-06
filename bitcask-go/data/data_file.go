package data

import "bitcask-go/fio"

const DataFileNameSuffix = ".data"

type DataFile struct {
	FileId    uint32        //文件Id
	WriteOff  int64         //文件写到了那个位置
	IoManager fio.IOManager //文件IO读写管理
}

// OpenDataFile 打开新的数据文件
func OpenDataFile(dirPath string, fileId uint32) (*DataFile, error) {
	return nil, nil
}

// Sync 同步当前活跃文件
func (df *DataFile) Sync() error {
	return nil
}

func (df *DataFile) Write(buf []byte) error {
	return nil
}

func (df *DataFile) ReadLogRecord(offset int64) (*LogRecord, int64, error) {
	return nil, 0, nil
}
