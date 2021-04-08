package dbmod

import (
	"encoding/json"
	"fmt"

	"log"
)

// RecDir 记录当前记录的操作位置
type RecDir struct {
	RecNo   int    `json:"sn" `   //记录流水号
	WriteID int    `json:"wid" `  //写的位置
	ReadID  int    `json:"rid" `  //读的位置
	Flag    string `json:"flag" ` //是否循环覆盖写
}

func InitRecDir() (err error) {
	for i := 0; i < MAXRECDIRS; i++ {
		rd := RecDir{}
		key := fmt.Sprintf("RecDirTB%02d", i+1)
		bv, err := json.Marshal(rd)
		if err != nil {
			log.Println("UpdateDirs Marshal Error:", err)
			return err
		}
		err = PutData(key, bv)
		if err != nil {
			return err
		}
	}
	return err
}

// UpdateDirs 更新目录
func (rd *RecDir) UpdateDirs(areaID RecArea) error {
	key := fmt.Sprintf("RecDirTB%02d", areaID)
	bv, err := json.Marshal(rd)
	if err != nil {
		log.Println("UpdateDirs Marshal Error:", err)
		return err
	}
	err = PutData(key, bv)
	return err
}

// LoadDirs 加载(读取)目录
func (rd *RecDir) LoadDirs(areaID RecArea) error {
	key := fmt.Sprintf("RecDirTB%02d", areaID)
	bv, err := GetData(key)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bv, rd)
	if err != nil {
		log.Println("LoadDirs Unmarshal Error:", err)
	}
	return err
}
