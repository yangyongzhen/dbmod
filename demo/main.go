package main

import (
	"dbmod"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

// RequestQrcode模拟写入的数据内容json
type RequestQrcode struct {
	ChanNO   string `json:"chanNO"`
	TermID   string `json:"termID"`
	Qrcode   string `json:"qrcode"`
	Money    uint32 `json:"money"`
	Recsn    uint32 `json:"recsn"`
	Orderno  string `json:"orderno"`
	Dealtime string `json:"dealtime"`
}

func saveREC(sn int) {
	defer wg.Done()
	//t := time.Now()
	//fmt.Println("ENTER")
	data := RequestQrcode{}
	data.ChanNO = "YS_CHANaaa"
	data.TermID = "12345678"
	data.Recsn = uint32(sn)
	data.Qrcode = "6225882618789"
	data.Money = 1
	// 按队列顺序写记录
	id, err := recApi.SaveRec(dbmod.RecArea01, data, 0x01)
	if err != nil {
		fmt.Printf("SaveRec error,%s\n", err.Error())
	} else {
		//非必须，更新记录的接口测试
		data.Money = 3
		_, err = recApi.UpdateRec(dbmod.RecArea01, id, data, 0x03)
		if err != nil {
			fmt.Printf("UpdateRec error,%s\n", err.Error())
		}
	}
	//elapsed := time.Since(t)
	//fmt.Println("saveEXIT,elapsed=", elapsed)
}

var recApi dbmod.Recorder
var wg sync.WaitGroup

func main() {
	fmt.Println("test record and queue")
	c := make(chan os.Signal)

	runtime.GOMAXPROCS(3)              // 限制 CPU 使用数，避免过载
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
	runtime.SetCPUProfileRate(10)
	//打开leveadmin查看数据　http://127.0.0.1:9999/leveldb_admin/static/
	go http.ListenAndServe(":9999", nil)

	recApi = dbmod.NewRecAPI(false)

	//初始化存储区，只需第一次执行一次,后续只需每次启动后OpenRecAreas
	err := recApi.InitRecAreas()
	if err != nil {
		fmt.Printf("InitRecAreas error,%s\n", err.Error())
	}
	//每次必须先打开存储区
	err = recApi.OpenRecAreas()
	if err != nil {
		fmt.Printf("OpenRecAreas error,%s\n", err.Error())
	}
	fmt.Println("ENTER")
	//recApi.DeleteRec(1, 1000)
	wg.Add(1000)
	//模拟异步并发写入1000条记录
	t := time.Now()
	for i := 0; i < 1000; i++ {
		go saveREC(i)
	}
	//saveREC()

	wg.Wait()
	elapsed := time.Since(t)
	fmt.Println("saveEXIT,耗时=", elapsed)

	start := time.Now().UnixNano()
	//删除999条未传记录,或者recApi.DeleteRec(1, 999)
	for i := 0; i < 1; i++ {
		recApi.DeleteRec(dbmod.RecArea01, 1)
	}
	//recApi.DeleteRec(1, 1)
	//记录结束时间
	end := time.Now().UnixNano()

	//输出执行时间，单位为毫秒。
	fmt.Println((end - start) / 1000000)
	//fmt.Println("elapsed=", elapsed)
	//recApi.DeleteRec(1, 5)
	//获取未传记录数量
	num := recApi.GetNoUploadNum(dbmod.RecArea01)
	fmt.Printf("GetNoUploadNum:%d\n", num)

	//读取一条未传记录,记录存储队列的头部开始读取
	rec, err := recApi.ReadRecNotServer(dbmod.RecArea01, 1)
	if err != nil {
		fmt.Printf("ReadRecNotServer error,%s\n", err.Error())
	}
	fmt.Printf("rec:%#v\n", rec)
	//读取下一条未传记录
	rec, err = recApi.ReadRecNotServer(dbmod.RecArea01, 2)
	if err != nil {
		fmt.Printf("ReadRecNotServer error,%s\n", err.Error())
	}
	fmt.Printf("rec:%#v\n", rec)

	//读取最近的一笔写入记录,从队列的尾部顺序读取
	rec, err = recApi.ReadRecWriteNot(dbmod.RecArea01, 1)
	if err != nil {
		fmt.Printf("ReadRecWriteNot error,%s\n", err.Error())
	}
	fmt.Printf("rec:%#v\n", rec)
	//读取最近的第二笔写入记录
	rec, err = recApi.ReadRecWriteNot(dbmod.RecArea01, 2)
	if err != nil {
		fmt.Printf("ReadRecWriteNot error,%s\n", err.Error())
	}
	fmt.Printf("rec:%#v\n", rec)
	//按ID读取记录
	rec, err = recApi.ReadRecByID(dbmod.RecArea01, 1)
	if err != nil {
		fmt.Printf("ReadRecByID error,%s\n", err.Error())
	}
	fmt.Printf("rec:%#v\n", rec)
	rec, err = recApi.ReadRecByID(dbmod.RecArea01, 10)
	if err != nil {
		fmt.Printf("ReadRecByID error,%s\n", err.Error())
	}
	fmt.Printf("rec:%#v\n", rec)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c

	fmt.Printf("exits,%v\n", s)
}
