package bitcask_go

import "errors"

var (
	ErrKeyIsEmpty        = errors.New("key is empty")
	ErrIndexUpdateFailed = errors.New("idx update failed")
	ErrNotFoundPos       = errors.New("not found logrecordpos")
	ErrDatafileNotFound  = errors.New("datafile not found")
	ErrKeyNotFound       = errors.New("key not found")
	ErrDataDirCorrupted  = errors.New("the database dir maybe corrupted")
)
