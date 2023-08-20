package CSql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var m_db *sql.DB


func OpenDb(ip, port, dbName, username, password string) (error) {
	
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+ip+":"+port+")/"+dbName)
	m_db =db
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = m_db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}
	m_db.SetMaxOpenConns(200)
	m_db.SetMaxIdleConns(100)
	fmt.Println("Connected to MySQL database!")
	return nil
}

