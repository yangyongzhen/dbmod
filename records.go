package dbmod

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	//IsDebug 是否调试
	IsDebug       = true
	recDir        [MAXRECAREAS]RecDir
	lockSave      = sync.Mutex{}
	lockDel       = sync.Mutex{}
	once          sync.Once
	singleintance *Records
)

// Records ...
type Records struct {
	ID      int         `json:"id"`
	RecNo   int         `json:"sn"`
	RecType int         `json:"type" `
	RecTime string      `json:"time" `
	Data    interface{} `json:"data" `
	Ext     string      `json:"ext" `
	Res     string      `json:"res" `
}

// InitRecAreas 初始化记录存储区
func (rec Records) InitRecAreas() error {

	//　清空数据
	err := DelAllData()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	// 初始化目录表
	err = InitRecDir()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return err
}

// OpenRecAreas 打开记录存储区,每次开机，需要先打开一下
func (rec *Records) OpenRecAreas() (err error) {
	//加载RecDir
	for i := 0; i < MAXRECAREAS; i++ {
		log.Printf("LoadDirs %02d \n", i+1)
		err = recDir[i].LoadDirs(RecArea(i) + 1)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Printf("LoadDirs %02d ok!\n", i+1)
	}
	//log.Println(recDir)

	return err
}

func saveRecToLDB(areaID RecArea, rec *Records, wid int) (id int, err error) {
	t := time.Now()
	rec.RecTime = t.Format("20060102150405")
	rec.ID = wid
	key := fmt.Sprintf("Rec%02dTB|%d", areaID, wid)
	bv, err := json.Marshal(rec)
	if err != nil {
		log.Println("saveRecToLDB Marshal Error:", err)
		return id, err
	}
	err = PutData(key, bv)
	if err != nil {
		log.Println("saveRecToLDB PutData Error:", err)
		return id, err
	}
	return id, err
}

// SaveRec 保存记录
func (rec *Records) SaveRec(areaID RecArea, data interface{}, recType int) (id int, err error) {
	lockSave.Lock()
	defer lockSave.Unlock()
	//log.Printf("SaveRec,area=%02d \n", areaID)
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		err = fmt.Errorf("area id  %02d is not right,must between 1 and %02d", areaID, MAXRECAREAS)
		log.Println(err.Error())
		return
	}
	rec.RecNo = recDir[areaID-1].RecNo
	rec.Data = data
	rec.RecType = recType
	//记录是否存储满，判断
	if (recDir[areaID-1].WriteID + 1) > (MAXRECCOUNTS) {

		if recDir[areaID-1].ReadID == 0 {

			err = fmt.Errorf("rec area %02d is full", areaID)
			log.Println(err.Error())
			return
		}

		if (recDir[areaID-1].WriteID + 1 - (MAXRECCOUNTS)) == recDir[areaID-1].ReadID {
			err = fmt.Errorf("rec area %02d is full", areaID)
			log.Println(err.Error())
			return
		}

		//保存记录

		recDir[areaID-1].RecNo++
		recDir[areaID-1].WriteID = 1
		recDir[areaID-1].Flag = "1"
		id = 1
		_, err = saveRecToLDB(areaID, rec, id)
		if err != nil {
			log.Println("saveRecToLDB Error:", err)
			return
		}
		err = recDir[areaID-1].UpdateDirs(areaID)
		if err != nil {
			log.Println("SaveRec UpdateDirs Error:", err)
			return
		}
		//log.Printf("SaveRec,area=%02d ok!\n", areaID)
		return id, err
	}

	if recDir[areaID-1].Flag == "1" {
		//记录是否满判断
		if (recDir[areaID-1].WriteID + 1) == recDir[areaID-1].ReadID {
			err = fmt.Errorf("rec area %02d is full", areaID)
			log.Println(err.Error())
			return
		}
		rec.RecNo += 1
		id = recDir[areaID-1].WriteID + 1
		_, err = saveRecToLDB(areaID, rec, id)
		if err != nil {
			log.Println("saveRecToLDB Error:", err)
			return
		}
		recDir[areaID-1].RecNo++
		recDir[areaID-1].WriteID = id
		err = recDir[areaID-1].UpdateDirs(areaID)
		if err != nil {
			log.Fatal(err.Error())
			return 0, err
		}
		//log.Printf("SaveRec,area=%02d ok!\n", areaID)
		return id, err

	}

	rec.RecNo += 1
	id = recDir[areaID-1].WriteID + 1
	_, err = saveRecToLDB(areaID, rec, id)
	if err != nil {
		log.Println("saveRecToLDB Error:", err)
		return
	}
	recDir[areaID-1].RecNo++
	recDir[areaID-1].WriteID = id
	err = recDir[areaID-1].UpdateDirs(areaID)
	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}
	//log.Printf("SaveRec,area=%02d ok!\n", areaID)
	return id, err

}

// UpdateRec 更新记录
func (rec *Records) UpdateRec(areaID RecArea, recID int, data interface{}, recType int) (id int, err error) {
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		err = fmt.Errorf("area id  %02d is not right,must between 1 and %02d", areaID, MAXRECAREAS)
		log.Println(err.Error())
		return
	}
	rec.Data = data
	rec.RecType = recType
	id, err = saveRecToLDB(areaID, rec, recID)
	return id, err
}

// DeleteRec 删除记录（并不是真正删除表里记录，而是清除该记录的上传标记）
// areaID:记录区 num:删除的数量
func (rec Records) DeleteRec(areaID RecArea, num int) (err error) {

	lockDel.Lock()
	defer lockDel.Unlock()
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		err = errors.New("area id is not right")
		log.Fatal(err.Error())
		return
	}

	id := recDir[areaID-1].ReadID

	//如果写的位置等于读的位置，说明记录已上传完，没有要删除的了
	if recDir[areaID-1].WriteID == recDir[areaID-1].ReadID {
		return
	}

	//如果要删除的数量大于了最大的记录数
	if (id + num) > MAXRECCOUNTS {
		if (id + num - MAXRECCOUNTS) > recDir[areaID-1].WriteID {
			recDir[areaID-1].ReadID = recDir[areaID-1].WriteID
			err = recDir[areaID-1].UpdateDirs(areaID)
			if err != nil {
				log.Fatal(err.Error())
				return err
			}
			return
		}
		//更新读指针（读的位置）
		recDir[areaID-1].ReadID = id + num - MAXRECCOUNTS
		err = recDir[areaID-1].UpdateDirs(areaID)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		return
	}

	//如果当前写的位置大于读的位置
	if recDir[areaID-1].WriteID > recDir[areaID-1].ReadID {
		if id+num > recDir[areaID-1].WriteID {
			//更新读指针（读的位置）
			recDir[areaID-1].ReadID = recDir[areaID-1].WriteID
			err = recDir[areaID-1].UpdateDirs(areaID)
			if err != nil {
				log.Fatal(err.Error())
				return err
			}
			return
		}
	}

	//更新读指针（读的位置）
	recDir[areaID-1].ReadID = id + num
	err = recDir[areaID-1].UpdateDirs(areaID)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return
}

//GetNoUploadNum 获取未上传记录数量
func (rec Records) GetNoUploadNum(areaID RecArea) int {

	num := 0
	if recDir[areaID-1].WriteID == recDir[areaID-1].ReadID {
		num = 0
		return num
	}
	if recDir[areaID-1].Flag != "1" {
		num = int(recDir[areaID-1].WriteID - recDir[areaID-1].ReadID)
	} else {
		if recDir[areaID-1].WriteID > recDir[areaID-1].ReadID {
			num = int(recDir[areaID-1].WriteID - recDir[areaID-1].ReadID)
		} else {
			num = int(MAXRECCOUNTS - recDir[areaID-1].ReadID + recDir[areaID-1].WriteID)
		}
	}
	return num
}

// ReadRecByID 按数据库ID读取记录
func (rec Records) ReadRecByID(areaID RecArea, rid int) (p *Records, err error) {
	var rec1 Records
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		err = errors.New("area id is not right")
		log.Fatal(err.Error())
		return
	}
	key := fmt.Sprintf("Rec%02dTB|%d", areaID, rid)
	bv, err := GetData(key)
	err = json.Unmarshal(bv, &rec1)
	if err != nil {
		log.Println("ReadRecByID Unmarshal Error:", key, err)
	}
	return &rec1, nil
}

//ReadRecNotServer 读取未上传的记录数据，顺序读取第SN条未上传的记录
//sn取值 1-到-->未上传记录数目
func (rec Records) ReadRecNotServer(areaID RecArea, sn int) (p *Records, err error) {
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		err = errors.New("area id is not right")
		log.Fatal(err.Error())
		return
	}
	id := recDir[areaID-1].ReadID
	//fmt.Printf("id=%d\n", id)
	if (int(id) + sn) > MAXRECCOUNTS {
		if int(id)+sn-MAXRECCOUNTS > int(recDir[areaID-1].WriteID) {
			return nil, errors.New("no records")
		}
		p, err = rec.ReadRecByID(areaID, int(id)+sn-MAXRECCOUNTS)
	} else {
		if recDir[areaID-1].ReadID <= recDir[areaID-1].WriteID {
			if (int(id) + sn) > int(recDir[areaID-1].WriteID) {
				return nil, errors.New("no records")
			}
		}
		p, err = rec.ReadRecByID(areaID, int(recDir[areaID-1].ReadID)+sn)

	}
	return p, err
}

// ReadRecWriteNot 倒数读取第SN条写入的记录
//读取一条记录  倒数读取第SN条写入的记录
func (rec Records) ReadRecWriteNot(areaID RecArea, sn int) (p *Records, err error) {
	id := int(recDir[areaID-1].WriteID)
	if (id - sn) < 0 {
		if recDir[areaID-1].Flag == "1" {
			p, err = rec.ReadRecByID(areaID, MAXRECCOUNTS-(sn-id-1))
		} else {
			return nil, errors.New("no records")
		}
	} else {
		p, err = rec.ReadRecByID(areaID, (id - sn + 1))
	}
	return
}

// GetLastRecNO 获取最后一条记录流水号
func (rec Records) GetLastRecNO(areaID RecArea) int {
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		log.Println("area id is not right")
		return 0
	}
	id := recDir[areaID-1].RecNo
	return id
}

// GetReadWriteID 获取当前的读、写ID
func (rec Records) GetReadWriteID(areaID RecArea) (rid, wid int) {
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		log.Println("area id is not right")
		return 0, 0
	}
	rid = recDir[areaID-1].ReadID
	wid = recDir[areaID-1].WriteID
	return rid, wid
}

// NewRecords ...
func NewRecords(debug bool) *Records {
	IsDebug = debug
	if singleintance == nil {
		once.Do(func() {
			fmt.Println("Init singleintance Record operation ")
			singleintance = new(Records)
		})
	}

	return singleintance
}
