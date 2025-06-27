package Adp

import (
	SQL "adpGo/pkg/CSql"
	"fmt"
	"strconv"
)

func GetNewCusId() string {

	var iIdx int
	in := make(map[string]interface{})
	var listOut []interface{}
	var sError string
	in["skey"] = "CusId"

	SQL.QueryTb("Settings", in, &listOut, &sError)

	if len(listOut) < 1 {

		tmp := make(map[string]interface{})
		var listTmp []interface{}
		tmp["skey"] = "CustomerId"
		SQL.QueryTb("Settings", tmp, &listTmp, &sError)
		if len(listTmp) > 0 {
			str := listTmp[0].(map[string]interface{})["svalue"].(string)
			iIdx, _ = ParseCusID(str)

		} else {
			iIdx = 0
		}

	} else {

		row := listOut[0].(map[string]interface{})
		sValueStr := row["svalue"].(string)
		var err error
		iIdx, err = strconv.Atoi(sValueStr)
		if err != nil {
			iIdx = 0 // 或記錄錯誤
		}

	}
	iIdx += 1
	in["svalue"] = strconv.Itoa(iIdx)
	fmt.Print(in)
	SQL.InsertTb("Settings", in, &sError, true)

	return GenerateCusID(iIdx)
}

// GenerateID 產生像 EA01, EB02, ..., FA01 的編號
func GenerateCusID(n int) string {
	if n < 1 {
		return ""
	}
	n -= 1 // 0-based index

	group := n / 99
	num := n%99 + 1

	first := 'E' + rune(group/26)
	second := 'A' + rune(group%26)

	if first > 'Z' {
		return "OVERFLOW"
	}

	return fmt.Sprintf("%c%c%02d", first, second, num)
}

// ParseID 將像 EA99、FA01 的 ID 轉回整數序號（第幾個）
func ParseCusID(id string) (int, error) {
	if len(id) != 4 {
		return 0, fmt.Errorf("invalid ID length")
	}

	first := id[0]
	second := id[1]
	numberPart := id[2:]

	if first < 'A' || first > 'Z' || second < 'A' || second > 'Z' {
		return 0, fmt.Errorf("invalid letters in ID")
	}

	num, err := strconv.Atoi(numberPart)
	if err != nil || num < 1 || num > 99 {
		return 0, fmt.Errorf("invalid number in ID")
	}

	group := (int(first-'E') * 26) + int(second-'A')
	n := group*99 + num

	return n, nil
}
