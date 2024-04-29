package data

// LogRecordPos 数据内存索引、主要是描述记录在文件中的位置
type LogRecordPos struct {
	Fid    uint32 //文件ID、表示将数据存储在那个块中
	Offset int64  //记录的偏移量

}
