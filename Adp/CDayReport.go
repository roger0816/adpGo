package Adp

import(

CSQL "github.com/roger0816/adpGo/CSql"
C "github.com/roger0816/adpGo/Common"

)

func GetDayReport(Data map[string]interface{}, sError string) ( bool ,[]interface{} ) {

	reList := []interface{}{}

	orderData := []interface{}{}

	var bOk bool
	bOk = CSQL.QueryTb(C.SQL_TABLE.OrderData(), Data, &orderData, &sError)

	return bOk,reList
}
