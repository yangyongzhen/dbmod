# dbmod
Encapsulation of operation leveldB records, a Implementation of queue。

按队列顺序操作leveldb记录存储模块的封装。

封装完成后的leveldb的记录存储操作多么的简单，且是一个持久化的队列的实现，只需以下接口：

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


func main() {
  var recApi dbmod.Recorder
  //每次必须先打开存储区
	err := recApi.OpenRecAreas()
	if err != nil {
		fmt.Printf("OpenRecAreas error,%s\n", err.Error())
	}
  
  // 按队列顺序写入一条记录,data为interface{},会序列化为json存储
	id, err := recApi.SaveRec(dbmod.RecArea01, data, datatype)
	if err != nil {
		fmt.Printf("SaveRec error,%s\n", err.Error())
	}
  //按队列顺序读取一条记录
  rec, err := recApi.ReadRecNotServer(dbmod.RecArea01, 1)
	if err != nil {
		fmt.Printf("ReadRecNotServer error,%s\n", err.Error())
	}
	fmt.Printf("rec:%#v\n", rec)
  
  //按队列顺序删除一条记录(注:只上传更新标记)
  recApi.DeleteRec(dbmod.RecArea01, 1)
  
  //获取队列中未消费的记录数量
	num := recApi.GetNoUploadNum(dbmod.RecArea01)
	fmt.Printf("GetNoUploadNum:%d\n", num)
}
