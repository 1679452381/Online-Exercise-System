package test

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

// gin连接mysql
func TestMysql(t *testing.T) {
	Db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/online_exercise")
	if err != nil {
		t.Fatal(err.Error())
	}
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)

}
