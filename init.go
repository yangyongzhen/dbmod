package dbmod

import (
	levelAdmin "github.com/qjues/leveldb-admin"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"net/http"
)

var LDB *leveldb.DB

func init() {
	var err error
	LDB, err = leveldb.OpenFile("./data/recordsDB", nil)
	if err != nil {
		log.Fatal("LevelDB create recordsDB error:", err)
		return
	}
	//http://127.0.0.1:9999/leveldb_admin/static/
	levelAdmin.GetLevelAdmin().Register(LDB, "recordsDB").SetServerMux(http.DefaultServeMux).Start()
	return
}
