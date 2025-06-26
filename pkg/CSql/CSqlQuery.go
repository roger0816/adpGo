package CSql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var m_db *sql.DB

var m_writeDb *sql.DB

var LastTbTime = make(map[string]interface{})

func db() *sql.DB {
	return m_db
}

// to do 讀寫分離
func writeDb() *sql.DB {
	return m_db
}

func OpenDb(ip, port, dbName, username, password string) error {
	fmt.Println("db open 1003")
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+ip+":"+port+")/"+dbName)
	m_db = db
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

func CloseDb() {
	if m_db != nil {
		m_db.Close()
	}
}

func Init(listTableName []string, updateTimeKey string) {

	for _, v := range listTableName {
		in := make(map[string]interface{})
		in["DESC"] = updateTimeKey
		in["LIMIT"] = "1"

		listOut := []interface{}{}
		var sError string
		bOk := QueryTb(v, in, &listOut, &sError)

		if bOk && len(listOut) > 0 {
			if record, ok := listOut[0].(map[string]interface{}); ok {
				if dateTime, exists := record[updateTimeKey]; exists {
					LastTbTime[v] = dateTime
				}
			}
		}
	}

	fmt.Printf("lastTbTime: %+v\n", LastTbTime)
}

func QueryTb(tableName string, conditions map[string]interface{}, listOut *[]interface{}, sError *string) bool {
	return BaseQuery(tableName, conditions, listOut, sError, false)
}

func BaseQuery(tableName string, conditions map[string]interface{}, listOut *[]interface{}, sError *string, UseMainDb bool) bool {
	*listOut = nil

	for k, v := range conditions {
		if v == nil {
			delete(conditions, k)
		}
	}

	query := "SELECT * FROM " + tableName
	args := make([]interface{}, 0)
	whereClauses := make([]string, 0)
	var sSub string = ""
	var sOrderBy string = ""
	var sLimit string = ""

	for key, value := range conditions {
		if strings.ToUpper(key) == "ASC" {
			sOrderBy = " ORDER BY " + value.(string)
			continue
		}
		if strings.ToUpper(key) == "DESC" {
			sOrderBy += " ORDER BY " + value.(string) + " DESC"
			continue
		}
		if strings.ToUpper(key) == "LIMIT" {
			sLimit += " LIMIT " + value.(string)
			continue
		}

		tmpShift := "=?"

		if len(strings.Split(key, " ")) > 1 {
			tmpShift = "?"
		}

		whereClauses = append(whereClauses, key+tmpShift)

		switch v := value.(type) {
		case bool:
			if v {
				args = append(args, "1")
			} else {
				args = append(args, "0")
			}
		default:
			args = append(args, value)
		}

	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if sOrderBy != "" {
		sSub += sOrderBy
	}

	if sLimit != "" {
		sSub += sLimit
	}

	//	fmt.Printf("sql query : %s \n", query+sSub)

	var stmt *sql.Stmt
	var err error
	if UseMainDb {
		stmt, err = writeDb().Prepare(query + sSub)
	} else {
		stmt, err = db().Prepare(query + sSub)
	}

	if err != nil {
		log.Println("Error preparing query:", err)
		*sError = err.Error()
		return false
	}
	defer stmt.Close()

	fmt.Printf("sql query : %s \n", interpolateQuery(query+sSub, args))
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

			switch v := (*rowData[i].(*interface{})).(type) {
			case []byte:
				data[name] = string(v)
			default:
				data[name] = v
			}

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

func QueryCount(sTableName string, conditions map[string]interface{}, useMainDb bool) int {
	iRe := 0

	var sSub string
	var listTmp []interface{}
	var iCount int

	query := "SELECT COUNT(*) FROM " + sTableName

	for key, value := range conditions {
		if iCount == 0 {
			sSub += "  WHERE "
		} else {
			sSub += " AND "
		}
		if strings.Contains(key, " ") {
			sSub += key + " ?" // 处理自带 >= <= 或 like
		} else {
			sSub += key + " = ? "
		}
		iCount++
		listTmp = append(listTmp, value)
	}

	var rows *sql.Rows
	var err error

	if useMainDb {
		rows, err = writeDb().Query(query+sSub, listTmp...)

	} else {
		rows, err = db().Query(query+sSub, listTmp...)

	}
	if err != nil {
		//log.Fatal(err)

	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&iRe)
		if err != nil {
			//	log.Fatal(err)
		}
	}

	return iRe
}

func UpdateTb(sTableName string, conditions, data map[string]interface{}, sError *string) bool {
	sCmd := "UPDATE " + sTableName + " SET "

	sDateTime := CurrentTime()

	for k, v := range conditions {
		if v == nil {
			delete(conditions, k)
		}
	}

	for k, v := range data {
		if v == nil {
			delete(data, k)
		}
	}

	if len(data) < 1 {
		return false
	}

	if sTableName != "LastUpdateTime" {
		data["UpdateTime"] = sDateTime
	}
	var listKey []string
	for key := range data {
		listKey = append(listKey, key)
	}

	bFirst := true
	for _, key := range listKey {
		v := data[key]

		if key == "Sid" || strings.TrimSpace(fmt.Sprint(v)) == "" || v == nil {
			continue
		}

		if !bFirst {
			sCmd += ", "
		}
		sCmd += key + " = ? "

		bFirst = false
	}

	var sSub string
	var tmp []string
	for key := range conditions {
		tmp = append(tmp, key)
	}

	for i, key := range tmp {
		if i == 0 {
			sSub += " WHERE "
		} else {
			sSub += " AND "
		}

		v := conditions[key]

		if boolType, ok := v.(bool); ok {
			if boolType {
				v = 1
			} else {
				v = 0
			}
		}

		if strType, ok := v.(string); ok {
			sSub += key + "='" + strType + "' "
		} else {
			sSub += key + "=" + fmt.Sprint(v) + " "
		}
	}

	query := sCmd + sSub

	stmt, err := writeDb().Prepare(query)
	if err != nil {
		//log.Fatal(err)
		fmt.Print(err)
	}
	defer stmt.Close()

	args := make([]interface{}, 0)
	for _, key := range listKey {
		v := data[key]

		if key == "Sid" || strings.TrimSpace(fmt.Sprint(v)) == "" || v == nil {
			continue
		}

		args = append(args, v)
	}

	_, execErr := stmt.Exec(args...)
	if execErr != nil {
		fmt.Println(execErr)
	}
	setLastUpdateTime(sTableName, sDateTime)

	return true
}
func BatchUpdateTb(sTableName string, conditionsList, dataList []map[string]interface{}, sError *string) bool {
	if len(dataList) != len(conditionsList) {
		*sError = "The length of conditionsList and dataList must be equal"
		return false
	}

	sDateTime := CurrentTime()

	baseCmd := "UPDATE " + sTableName + " SET "

	transaction, err := writeDb().Begin()
	if err != nil {
		//log.Fatal(err)
		return false
	}

	for idx, data := range dataList {
		conditions := conditionsList[idx]

		if sTableName != "LastUpdateTime" {
			data["UpdateTime"] = sDateTime
		}

		var listKey []string
		for key := range data {
			listKey = append(listKey, key)
		}

		sCmd := ""
		bFirst := true
		for _, key := range listKey {
			v := data[key]

			if key == "Sid" || strings.TrimSpace(fmt.Sprint(v)) == "" || v == nil {
				continue
			}

			if !bFirst {
				sCmd += ", "
			}
			sCmd += key + " = ? "

			bFirst = false
		}

		var sSub string
		var tmp []string
		for key := range conditions {
			tmp = append(tmp, key)
		}

		for i, key := range tmp {
			if i == 0 {
				sSub += " WHERE "
			} else {
				sSub += " AND "
			}

			v := conditions[key]

			if boolType, ok := v.(bool); ok {
				if boolType {
					v = 1
				} else {
					v = 0
				}
			}

			if strType, ok := v.(string); ok {
				sSub += key + "='" + strType + "' "
			} else {
				sSub += key + "=" + fmt.Sprint(v) + " "
			}
		}

		query := baseCmd + sCmd + sSub
		stmt, err := transaction.Prepare(query)
		if err != nil {
			//log.Fatal(err)
			return false
		}

		args := make([]interface{}, 0)
		for _, key := range listKey {
			v := data[key]

			if key == "Sid" || strings.TrimSpace(fmt.Sprint(v)) == "" || v == nil {
				continue
			}

			args = append(args, v)
		}

		_, execErr := stmt.Exec(args...)
		if execErr != nil {
			//log.Fatal(execErr)
			transaction.Rollback()
			return false
		}
		stmt.Close()
	}

	err = transaction.Commit()
	if err != nil {
		//log.Fatal(err)
		return false
	}

	setLastUpdateTime(sTableName, sDateTime)

	return true
}

func GenerateId(owner string, date string) (int, error) {
	var id int

	// 當前時間作為 UpdateTime
	updateTime := CurrentTime()

	// 開始事務
	db := writeDb()
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// 確保在函數結束時根據情況提交或回滾事務
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 鎖定特定 owner 和日期的行
	err = tx.QueryRow(`
        SELECT Id 
        FROM MakeId 
        WHERE Owner = ? AND Date = ? 
        FOR UPDATE`, owner, date).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果沒有記錄，插入初始值 1
			id = 1
			_, err = tx.Exec(`
                INSERT INTO MakeId (Id, Date, Owner, UpdateTime) 
                VALUES (?, ?, ?, ?)`, id, date, owner, updateTime)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	} else {
		// 如果記錄存在，遞增 Id
		id += 1
		_, err = tx.Exec(`
            UPDATE MakeId 
            SET Id = ?, UpdateTime = ? 
            WHERE Owner = ? AND Date = ?`, id, updateTime, owner, date)
		if err != nil {
			return 0, err
		}
	}

	// 返回新的 Id
	return id, nil
}

func LockInsertTb(sTableName string, input map[string]interface{}, sError *string, bOrReplace bool) (bool, int64, map[string]interface{}) {
	data := input

	sDateTime := CurrentTime()

	if sTableName != "LastUpdateTime" {
		data["UpdateTime"] = sDateTime // 使用 "updateTime"，確保大小寫正確
	}

	db := writeDb()

	// 开始一个新的事务
	tx, err := db.Begin()
	if err != nil {
		*sError = "Failed to start transaction: " + err.Error()
		return false, 0, nil
	}

	_, err = tx.Exec("LOCK TABLES " + sTableName + " WRITE")
	if err != nil {
		*sError = "Failed to lock table: " + err.Error()
		tx.Rollback()
		return false, 0, nil
	}

	// ...[构建sCmd, 准备args]...

	var listKey []string
	for k, v := range data {
		// 如果v是string類型且不是空字符串，或者v不是string且不是nil
		if str, ok := v.(string); ok {
			if str != "" {
				listKey = append(listKey, k)
			}
		} else if v != nil {
			listKey = append(listKey, k)
		}
	}

	operation := "INSERT"
	if bOrReplace {
		operation = "REPLACE"
	}

	var tmpKey, tmpValue string
	for i, sKey := range listKey {
		if i != 0 {
			tmpKey += ","
			tmpValue += ","
		}
		tmpKey += sKey
		tmpValue += "?"
	}

	sCmd := fmt.Sprintf("%s INTO %s (%s) VALUES (%s);", operation, sTableName, tmpKey, tmpValue)

	fmt.Println("cmd:", sCmd)

	// 使用事务执行插入操作

	args := make([]interface{}, len(listKey))
	for i, sKey := range listKey {
		args[i] = data[sKey]
		if sTableName != "PicData" {
			fmt.Println(sKey, ":", data[sKey])
		}
	}
	// Capture the result after execution
	fmt.Printf("sql insert : %s \n", interpolateQuery(sCmd, args))
	result, err := tx.Exec(sCmd, args...)
	if err != nil {
		*sError = err.Error()
		return false, 0, nil
	}

	_, err = tx.Exec("UNLOCK TABLES;")
	if err != nil {
		*sError = "Failed to lock table: " + err.Error()
		tx.Rollback()
		return false, 0, nil
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		*sError = "Failed to commit transaction: " + err.Error()
		return false, 0, nil
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		*sError = "Failed to retrieve the last insert ID: " + err.Error()
		return false, 0, nil
	}

	// Query the inserted data using a new Query and scan into the map
	rows, err := writeDb().Query("SELECT * FROM "+sTableName+" WHERE Sid=?", lastInsertID)
	if err != nil {
		*sError = "Failed to query the inserted data: " + err.Error()
		return true, lastInsertID, nil
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}
	if rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			*sError = "Failed to scan the inserted data: " + err.Error()
			return true, lastInsertID, nil
		}
	}

	resultData := make(map[string]interface{})

	var tmp = make(map[string]interface{})
	tmp["Sid"] = lastInsertID

	var listOut = []interface{}{}
	var tmpErr string
	QueryTb(sTableName, tmp, &listOut, &tmpErr)

	if len(listOut) > 0 {
		interFaceToMap(listOut[0], &resultData)
	}

	return true, lastInsertID, resultData
}

func InsertTb(sTableName string, input map[string]interface{}, sError *string, bOrReplace bool) (bool, int64, map[string]interface{}) {
	data := input

	sDateTime := CurrentTime()

	if sTableName != "LastUpdateTime" {
		data["UpdateTime"] = sDateTime // 使用 "updateTime"，確保大小寫正確
	}
	var listKey []string
	for k, v := range data {
		// 如果v是string類型且不是空字符串，或者v不是string且不是nil
		if str, ok := v.(string); ok {
			if str != "" {
				listKey = append(listKey, k)
			}
		} else if v != nil {
			listKey = append(listKey, k)
		}
	}

	operation := "INSERT"
	if bOrReplace {
		operation = "REPLACE"
	}

	var tmpKey, tmpValue string
	for i, sKey := range listKey {
		if i != 0 {
			tmpKey += ","
			tmpValue += ","
		}
		tmpKey += sKey
		tmpValue += "?"
	}

	sCmd := fmt.Sprintf("%s INTO %s (%s) VALUES (%s);", operation, sTableName, tmpKey, tmpValue)

	fmt.Println("cmd:", sCmd)

	query, err := writeDb().Prepare(sCmd)
	if err != nil {
		*sError = "Failed to query the inserted data: " + err.Error()
		//log.Fatal(err)
	}

	args := make([]interface{}, len(listKey))
	for i, sKey := range listKey {
		args[i] = data[sKey]
		if sTableName != "PicData" {
			fmt.Println(sKey, ":", data[sKey])
		}
	}
	// Capture the result after execution
	fmt.Printf("sql insert : %s \n", interpolateQuery(sCmd, args))
	result, err := query.Exec(args...)
	if err != nil {
		*sError = err.Error()
		return false, 0, nil
	}

	defer query.Close()

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		*sError = "Failed to retrieve the last insert ID: " + err.Error()
		return false, 0, nil
	}

	// Query the inserted data using a new Query and scan into the map
	rows, err := writeDb().Query("SELECT * FROM "+sTableName+" WHERE Sid=?", lastInsertID)
	if err != nil {
		*sError = "Failed to query the inserted data: " + err.Error()
		return true, lastInsertID, nil
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}
	if rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			*sError = "Failed to scan the inserted data: " + err.Error()
			return true, lastInsertID, nil
		}
	}

	resultData := make(map[string]interface{})

	var tmp = make(map[string]interface{})
	tmp["Sid"] = lastInsertID

	var listOut = []interface{}{}
	var tmpErr string
	QueryTb(sTableName, tmp, &listOut, &tmpErr)

	if len(listOut) > 0 {
		interFaceToMap(listOut[0], &resultData)
	}

	return true, lastInsertID, resultData
}

func DelFromTb(sTableName string, conditions map[string]interface{}, sError *string) bool {
	var whereClauses []string
	var args []interface{}

	// Constructing WHERE clauses based on conditions
	for k, v := range conditions {
		whereClauses = append(whereClauses, k+" = ?")
		args = append(args, v)
	}
	whereStr := strings.Join(whereClauses, " AND ")

	// Preparing the DELETE statement
	sCmd := fmt.Sprintf("DELETE FROM %s WHERE %s", sTableName, whereStr)
	fmt.Println("del from tb cmd:", sCmd)

	query, err := writeDb().Prepare(sCmd)
	if err != nil {
		*sError = err.Error()
		log.Println("Error in preparing query:", err)
		return false
	}
	defer query.Close()

	// Executing the DELETE statement
	result, err := query.Exec(args...)
	if err != nil {
		*sError = err.Error()
		log.Println("Error in executing query:", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		*sError = err.Error()
		log.Println("Error in getting rows affected:", err)
		return false
	}

	if rowsAffected == 0 {
		log.Println("No data were deleted.")
		return false
		// Optionally, you can return false here if you consider no deletions as a failure.
	}

	// You can setTrigger here if needed
	// setTrigger(sTableName, sDateTime)

	return true
}

func CurrentTime() string {

	location, err := time.LoadLocation("Asia/Taipei") // Taipei is in the UTC+8 timezone
	if err != nil {
		panic(err)
	}

	re := time.Now().In(location).Format("20060102150405")

	return re
}

func setLastUpdateTime(tableName string, sDateTime string) {

	if tableName == "LastUpdateTime" {
		return
	}

	in := make(map[string]interface{})
	in["TableName"] = tableName
	in["Last"] = sDateTime
	var sError string
	InsertTb("LastUpdateTime", in, &sError, true)

}

func LastCustomerId(sClassSid, sClassId string, sError *string) (string, bool) {
	out := sClassId + "-EA00"

	queryStr := "SELECT Id FROM CustomerData WHERE Class=? ORDER BY Id DESC"
	row := writeDb().QueryRow(queryStr, sClassSid)
	var id string
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return out, true
		}
		*sError = err.Error()
		return out, false
	}

	return id, true
}

func getLastUpdateTime(tableName string) string {
	if tableName == "LastUpdateTime" {
		return ""
	}

	var re string

	in := make(map[string]interface{})
	listOut := []interface{}{}
	in["TableName"] = tableName
	var sError string
	QueryTb("LastUpdateTime", in, &listOut, &sError)

	if len(listOut) > 0 {
		if v, ok := listOut[0].(map[string]interface{}); ok {
			if updateTimeStr, ok := v["Last"].(string); ok {
				re = updateTimeStr
			}

		}
	}

	return re
}

func LastCustomerAddCostID() (string, bool) {

	location, err0 := time.LoadLocation("Asia/Taipei") // Taipei is in the UTC+8 timezone
	if err0 != nil {
		panic(err0)
	}

	sDate := time.Now().In(location).Format("20060102")

	var id string
	tmpDate := sDate[2:]
	defaultID := tmpDate + "-0000"
	query := fmt.Sprintf("SELECT * FROM CustomerCost WHERE OrderTime LIKE '%%%s%%' ORDER BY OrderId DESC;", sDate)

	row := db().QueryRow(query)
	err := row.Scan(&id) // Assuming OrderId is the first column, adjust accordingly if not.
	if err == sql.ErrNoRows {
		return defaultID, false
	}
	if err != nil {
		return "", false
	}

	return id, true
}

func interFaceToMap(in interface{}, out *map[string]interface{}) {
	bytes, _ := json.Marshal(in)
	json.Unmarshal(bytes, &out)
}

func interpolateQuery(query string, args []interface{}) string {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			query = strings.Replace(query, "?", "'"+v+"'", 1)
		case int:
			query = strings.Replace(query, "?", strconv.Itoa(v), 1)
		case float64:
			query = strings.Replace(query, "?", fmt.Sprintf("%f", v), 1)
		// ... 其他数据类型，根据需要添加
		default:
			query = strings.Replace(query, "?", fmt.Sprintf("%v", v), 1)
		}
	}
	return query
}

func LastOrderId(sDate string) (string, error) {
	tmpDate := sDate[2:]
	sId := tmpDate + "-A000"

	sCmd := fmt.Sprintf("SELECT Id FROM OrderData WHERE OrderDate='%s' AND Id!='' ORDER BY Id DESC LIMIT 1", sDate)

	fmt.Println("cmd: ", sCmd)

	row := db().QueryRow(sCmd)

	var id sql.NullString
	if err := row.Scan(&id); err != nil {
		return sId, nil
	}

	if id.Valid {
		sId = id.String
	}
	fmt.Printf("AAAA :LastOrderId : %v \n", sId)
	return sId, nil
}

func LastOrderName(ownerSid string, date string) (string, error) {
	cmd := fmt.Sprintf("SELECT Name FROM OrderData WHERE OrderDate='%s' AND Owner='%s' ORDER BY Id DESC", date, ownerSid)

	fmt.Println("cmd:", cmd)

	rows, err := db().Query(cmd)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var name string
	if rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return "", err
		}
	}

	return name, nil
}
