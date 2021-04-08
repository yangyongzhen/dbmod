package dbmod

import (
	"fmt"
	"log"
)

// PutData 存数据
func PutData(key string, data []byte) (err error) {
	if IsDebug {
		fmt.Printf("PutData,key=%s\n", key)
	}
	err = LDB.Put([]byte(key), data, nil)
	if err != nil {
		log.Println("lb Put Error:", err)
		return err
	}
	//log.Println("Put Data ok:", val)
	return err
}

// GetData 取数据
func GetData(key string) (data []byte, err error) {
	if IsDebug {
		fmt.Printf("GetData,key=%s\n", key)
	}
	data, err = LDB.Get([]byte(key), nil)
	if err != nil {
		log.Println("lb Get Error:", err)
		return data, err
	}
	return data, nil
}

func DelAllData() (err error) {
	iter := LDB.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		//value := iter.Value()
		LDB.Delete(key, nil)
	}
	iter.Release()
	err = iter.Error()
	return
}
