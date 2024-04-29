package fio

import "os"

// FileIO 标准系统文件IO
type FileIO struct {
	fd *os.File //文件描述符
}

func NewFileIOManager(name string) (*FileIO, error) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, DataFilePerm)
	if err != nil {
		return nil, err
	} else {
		return &FileIO{
			file,
		}, nil
	}
}

// Read 从文件中根据具体位置读取对应数据
func (fio *FileIO) Read(data []byte, offset int64) (int, error) {
	return fio.fd.ReadAt(data, offset)
}

// Write 写入数据到文件中
func (fio *FileIO) Write(data []byte) (int, error) {
	return fio.fd.Write(data)
}

// Sync 同步内存映射
func (fio *FileIO) Sync() error {
	return fio.fd.Sync()
}

// Close 关闭文件
func (fio *FileIO) Close() error {
	return fio.fd.Close()
}
