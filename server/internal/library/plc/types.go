// Package plc 西门子 S7 PLC 通信驱动
// 协议：S7 TCP（以太网），适用于 S7-200 SMART / S7-300 / S7-400 / S7-1200 / S7-1500
package plc

// Area PLC 存储区类型
type Area string

const (
	AreaDB Area = "DB" // 数据块
	AreaM  Area = "M"  // 位存储区（Merker）
	AreaI  Area = "I"  // 输入映像区
	AreaQ  Area = "Q"  // 输出映像区
	AreaV  Area = "V"  // 变量存储区（S7-200 SMART 专用，等同于 DB1）
)

// DataType PLC 数据类型
type DataType string

const (
	TypeBool   DataType = "Bool"
	TypeByte   DataType = "Byte"
	TypeWord   DataType = "Word"
	TypeDWord  DataType = "DWord"
	TypeInt    DataType = "Int"
	TypeDInt   DataType = "DInt"
	TypeReal   DataType = "Real"
	TypeString DataType = "String"
)

// ReadRequest 单个数据点读取请求
type ReadRequest struct {
	PointID    int      // 数据点 ID（用于结果回传）
	Field      string   // 字段标识
	Area       Area     // 存储区
	DBNumber   int      // DB 块号（Area=DB 时有效）
	ByteOffset int      // 字节偏移
	BitOffset  int      // 位偏移（TypeBool 时有效）
	DataType   DataType // 数据类型
	Scale      float64  // 换算系数
	OffsetVal  float64  // 换算偏移
	AlarmMin   *float64 // 报警下限
	AlarmMax   *float64 // 报警上限
	Unit       string   // 单位
}

// ReadResult 单个数据点读取结果
type ReadResult struct {
	PointID    int
	Field      string
	RawValue   string   // 原始值字符串
	EngValue   float64  // 工程值（换算后）
	Valid      bool     // 是否读取成功
	AlarmType  int      // 0=无报警 1=超上限 2=超下限
	AlarmMin   *float64
	AlarmMax   *float64
	Unit       string
}
