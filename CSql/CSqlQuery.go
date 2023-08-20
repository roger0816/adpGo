package CSql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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


func queryTb(db *sql.DB, tableName string, conditions map[string]interface{}, listOut *[]map[string]interface{}, sError *string) bool {
	*listOut = nil

	query := "SELECT * FROM " + tableName
	args := make([]interface{}, 0)
	whereClauses := make([]string, 0)

	for key, value := range conditions {
		if strings.ToUpper(key) == "ASC" {
			query += " ORDER BY " + value.(string)
			continue
		}
		if strings.ToUpper(key) == "DESC" {
			query += " ORDER BY " + value.(string) + " DESC"
			continue
		}
		if strings.ToUpper(key) == "LIMIT" {
			query += " LIMIT " + value.(string)
			continue
		}

		whereClauses = append(whereClauses, key+"=?")
		args = append(args, value)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Error preparing query:", err)
		*sError = err.Error()
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		log.Println("Error executing query:", err)
		*sError = err.Error()
		return false
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		log.Println("Error getting column names:", err)
		*sError = err.Error()
		return false
	}

	for rows.Next() {
		rowData := make([]interface{}, len(columnNames))
		for i := range columnNames {
			rowData[i] = new(interface{})
		}
		err := rows.Scan(rowData...)
		if err != nil {
			log.Println("Error scanning row:", err)
			*sError = err.Error()
			return false
		}

		data := make(map[string]interface{})
		for i, name := range columnNames {
			data[name] = *rowData[i].(*interface{})
		}
		*listOut = append(*listOut, data)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after iterating rows:", err)
		*sError = err.Error()
		return false
	}

	return true
}
