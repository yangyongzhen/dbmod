package dbmod

// 配置项
const (
	// MAXRECDIRS 最大记录目录数量
	//(一个记录目录对应控制一个记录表,它记录了记录表中的数据存储和读取的位置)
	MAXRECDIRS = (3)
	// MAXRECAREAS 最大记录区数量 10个（即记录表的个数，必须跟记录目录数量保持一致）
	MAXRECAREAS = MAXRECDIRS
	// MAXRECCOUNTS 最大记录条数（即一个表中允许存储的最大记录数，
	// 记录存满后且上传标记已清除后，则从头开始存储覆盖，存储一条，覆盖一条）
	MAXRECCOUNTS = (11000)
)

type RecArea int

//枚举,记录区定义(一个记录区对应一个表)
const (
	//RecArea01 记录区1
	RecArea01 RecArea = iota + 1
	//RecArea02 记录区2
	RecArea02
	//RecArea03 记录区3
	RecArea03
	//RecArea04 记录区4
	RecArea04
	//RecArea05 记录区5
	RecArea05
	//RecArea06 记录区6
	RecArea06
	//RecArea07 记录区7
	RecArea07
	//RecArea08 记录区8
	RecArea08
	//RecArea09 记录区9
	RecArea09
	//RecArea10 记录区10
	RecArea10
)

// Recorder 操作记录的接口声明
type Recorder interface {
	// 初始化记录区(会清空所有数据!)
	InitRecAreas() error
	// 打开记录区(开机必须先打开一次)
	OpenRecAreas() (err error)
	// 保存记录 areaID,存储区ID,取值从1--至--MAXRECAREAS(相当于一个表)
	SaveRec(areaID RecArea, data interface{}, recType int) (id int, err error)
	// 更新记录
	UpdateRec(areaID RecArea, recID int, data interface{}, recType int) (id int, err error)
	// 删除记录
	DeleteRec(areaID RecArea, num int) (err error)
	// 获取未上传记录数量
	GetNoUploadNum(areaID RecArea) int
	// 按数据库中的ID读取一条记录
	ReadRecByID(areaID RecArea, id int) (p *Records, err error)
	// 顺序读取未上传的记录
	ReadRecNotServer(areaID RecArea, sn int) (p *Records, err error)
	// 倒数读取记录（如sn=1代表最后一次写入的记录）
	ReadRecWriteNot(areaID RecArea, sn int) (p *Records, err error)
	// 最后一条记录流水
	GetLastRecNO(areaID RecArea) int
	// 获取当前的读、写ID
	GetReadWriteID(areaID RecArea) (rid, wid int)
}

// RecAPI 操作接口的类
type RecAPI struct {
	Recorder
}

// NewRecAPI 初始化操作接口
func NewRecAPI(debug bool) RecAPI {
	return RecAPI{Recorder: NewRecords(debug)}
}
