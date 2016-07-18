package g

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

// TODO 草草的写了一个db连接池,优化下
var (
	dbLock sync.RWMutex
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = makeDbConn()
	if DB == nil || err != nil {
		log.Fatalln("g.InitDB, get db conn fail", err)
	}

	log.Println("g.InitDB ok")
}

func GetDbConn(connName string) (c *sql.DB, e error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	return DB, nil
}

// 创建一个新的mysql连接
func makeDbConn() (conn *sql.DB, err error) {
	conn, err = sql.Open("mysql", Config().DB.Dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(Config().DB.MaxOpen)
	conn.SetMaxIdleConns(Config().DB.MaxIdle)
	err = conn.Ping()

	return conn, err
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
