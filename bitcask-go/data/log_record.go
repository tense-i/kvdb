package data

type LogRecordType = byte

const (
	LogRecordNormal = iota
	LogRecordDeleted
)

// LogRecordPos 数据内存索引、主要是描述记录在文件中的位置
type LogRecordPos struct {
	Fid    uint32 //文件ID、表示将数据存储在那个块中
	Offset int64  //记录的偏移量
}

// LogRecord 写入到数据库文件的记录
type LogRecord struct {
	Key   []byte
	Value []byte

	//文件的类型、枚举类型
	Type LogRecordType
}

// EncodeLogRecord 对 logRecord  进行编码，返回字节数组和长度
func EncodeLogRecord(logRecord *LogRecord) ([]byte, int64) {
	return nil, 0
}
