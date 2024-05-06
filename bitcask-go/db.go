package bitcask_go

import (
	"bitcask-go/data"
	"bitcask-go/index"
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// DB bitcask存储引擎实例
type DB struct {
	option         Options //数据库选项
	mu             *sync.RWMutex
	activeFile     *data.DataFile            //当前的活跃文件、可以用于写入
	index          index.Indexer             //内存索引
	oldactiveFiles map[uint32]*data.DataFile //旧的文件对象、只能用于读文件
	fileIds        []int                     //有序的文件ID列表、只能在加载索引时使用
}

func Open(options Options) (*DB, error) {
	err := checkOptions(options)
	if err != nil {
		return nil, err
	}
	//判断目录是否存在、不存在则创建
	if _, err := os.Stat(options.DirPath); os.IsNotExist(err) {
		err := os.MkdirAll(options.DirPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	//初始化DB实例结构体

	db := &DB{
		option:         options,
		mu:             new(sync.RWMutex),
		activeFile:     nil,
		index:          index.NewIndexer(options.indexType),
		oldactiveFiles: make(map[uint32]*data.DataFile),
	}

	//加载数据文件
	err = db.localDatafile()
	if err != nil {
		return nil, err
	}

	//从数据文件中加载索引
	//加载索引文件
	err = db.localIndexDatafile()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func checkOptions(options Options) error {
	if options.DirPath == "" {
		return errors.New("文件夹不存在")
	}

	if options.DataFileMaxSize <= 0 {
		return errors.New("file size error")
	}
	return nil
}

func (db *DB) Put(key []byte, value []byte) error {
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}

	logRecord := &data.LogRecord{
		Value: value,
		Key:   key,
		Type:  data.LogRecordNormal,
	}

	//追加到当前活跃文件中
	pos, err := db.appendRecord(logRecord)
	if err != nil {
		return err
	}

	//更新内存索引
	if !db.index.Put(key, pos) {
		return ErrIndexUpdateFailed
	}
	return nil
}

func (db *DB) appendRecord(rec *data.LogRecord) (*data.LogRecordPos, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	//当前活跃文件是否存在
	if db.activeFile == nil {
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	//写入数据编码
	enRecord, size := data.EncodeLogRecord(rec)

	//如果写入的数据已经到达活跃文件的阈值、则关闭当前的活跃文件、打开新的活跃文件
	if db.activeFile.WriteOff+size > db.option.DataFileMaxSize {

		//先持久化当前的数据文件、确保数据已经写进磁盘、以后不会在对该文件执行写操作
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}
		//将当前的活跃文件放到oldactivedataFile中
		db.oldactiveFiles[db.activeFile.FileId] = db.activeFile

		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	writeOff := db.activeFile.WriteOff
	if err := db.activeFile.Write(enRecord); err != nil {
		return nil, err
	}

	//根据用户配置、是否写完立即同步
	if db.option.SyncWrite {
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}
	}

	pos := &data.LogRecordPos{
		Fid:    db.activeFile.FileId,
		Offset: writeOff,
	}

	return pos, nil
}

// setActiveDataFile 设置当前活跃文件
// 在访问此方法前必须持有互斥锁
func (db *DB) setActiveDataFile() error {
	var initialFileId uint32 = 0
	if db.activeFile != nil {
		initialFileId = db.activeFile.FileId + 1
	}
	datafile, err := data.OpenDataFile(db.option.DirPath, initialFileId)
	if err != nil {
		return err
	}
	db.activeFile = datafile
	return nil
}

// Get 根据key获取数据
func (db *DB) Get(key []byte) ([]byte, error) {
	//加读锁
	db.mu.RLock()
	defer db.mu.RUnlock()

	if len(key) == 0 {
		return nil, ErrKeyIsEmpty
	}
	//从内存中取出key对应的索引信息
	logRecordPos := db.index.Get(key)
	if logRecordPos == nil {
		return nil, ErrNotFoundPos
	}

	//根据文件id找到对应的数据文件
	var datafile *data.DataFile
	if db.activeFile.FileId == logRecordPos.Fid {
		datafile = db.activeFile
	} else {
		datafile = db.oldactiveFiles[logRecordPos.Fid]
	}

	//数据文件不存在
	if datafile == nil {
		return nil, ErrDatafileNotFound
	}

	//数据文件存在、根据offset读取数据
	logRecord, _, err := datafile.ReadLogRecord(logRecordPos.Offset)
	if err != nil {
		return nil, err
	}
	if logRecord.Type == data.LogRecordDeleted {
		return nil, ErrKeyNotFound
	}

	return logRecord.Value, nil

}

func (db *DB) localDatafile() error {
	dirEntries, err := os.ReadDir(db.option.DirPath)
	if err != nil {
		return err
	}
	var fileIds []int

	//遍历目录中与data为后缀的文件
	for _, entry := range dirEntries {
		if strings.HasSuffix(entry.Name(), data.DataFileNameSuffix) {
			splitNames := strings.Split(entry.Name(), ".")
			fileId, err := strconv.Atoi(splitNames[0])
			if err != nil {
				return ErrDataDirCorrupted
			}
			fileIds = append(fileIds, fileId)
		}
	}
	//对文件ID进行排序
	sort.Ints(fileIds)
	db.fileIds = fileIds
	//依次对每个文件执行打开操作
	for i, fid := range fileIds {
		datafile, err := data.OpenDataFile(db.option.DirPath, uint32(fid))
		if err != nil {
			return err
		}
		//最后一个文件为活跃文件
		if i == len(fileIds)-1 {
			db.activeFile = datafile
		} else {
			db.oldactiveFiles[uint32(fid)] = datafile
		}
	}
	return nil
}

// 从数据文件中加载索引
// 遍历文件中的每条记录、并更新到内存索引中
func (db *DB) localIndexDatafile() error {
	if len(db.fileIds) == 0 {
		return nil
	}

	//遍历数据文件
	for i, fid := range db.fileIds {
		var fileId = uint32(fid)
		var datafile *data.DataFile

		//当前文件为活跃文件
		if fileId == db.activeFile.FileId {
			datafile = db.activeFile
		} else {
			datafile = db.oldactiveFiles[fileId]
		}

		var offset int64 = 0
		//读取记录
		for {
			logRecord, size, err := datafile.ReadLogRecord(offset)

			//区分正常读取完EOF与错误读取
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			//构造索引结构、保存
			logRecordPos := &data.LogRecordPos{
				Fid:    fileId,
				Offset: offset,
			}

			//为待删除记录
			if logRecord.Type == data.LogRecordDeleted {
				db.index.Delete(logRecord.Key)
			} else {
				db.index.Put(logRecord.Key, logRecordPos)
			}

			offset += size
		}
		//记录活跃文件的最后写入位置
		if i == len(db.fileIds)-1 {
			db.activeFile.WriteOff = offset
		}
	}
	return nil
}
