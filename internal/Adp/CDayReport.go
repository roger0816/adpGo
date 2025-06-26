package Adp

import (
	C "adpGo/common"
	CSQL "adpGo/pkg/CSql"
)

func GetDayReport(Data map[string]interface{}, sError string) (bool, []interface{}) {

	reList := []interface{}{}

	orderData := []interface{}{}

	var bOk bool
	bOk = CSQL.QueryTb(C.SQL_TABLE.OrderData(), Data, &orderData, &sError)

	return bOk, reList
}
