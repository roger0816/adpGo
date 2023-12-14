package Adp

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	CSQL "github.com/roger0816/adpGo/CSql"
	C "github.com/roger0816/adpGo/Common"
	NETWORK "github.com/roger0816/adpGo/RpkNetwork"
)

func DoOrder(oriData NETWORK.CData, reData *C.VariantMap, reList *[]interface{}) (bOk bool, sOkMsg string, sError string) {

	var tmpIn C.VariantMap = make(map[string]interface{})

	//var tmpMap C.VariantMap = make(map[string]interface{})
	// var tmpList = []interface{}{}
	var Data C.VariantMap = oriData.Data

	var order, current C.OrderData

	// if sid, exists := Data["Sid"]; exists {
	// 	sidStr, isString := sid.(string)
	// 	if isString {
	// 		sidInt, err := strconv.Atoi(sidStr)
	// 		if err != nil {
	// 		//	return fmt.Errorf("error converting Sid from string to int: %v", err)
	// 		}
	// 		Data["Sid"] = sidInt
	// 	}
	// }

	C.MapToStruct(Data, &order)

	if order.Step == "0" {

		bOk, sError = orderStep0(&order)
		sOkMsg = "報價成功"
	} else {
		if order.Owner == "未分配" {
			current = order
		} else {
			tmpIn["Sid"] = C.Int64ToString(int64(order.Sid))
			var listOut []interface{}
			CSQL.QueryTb(C.SQL_TABLE.OrderData(), tmpIn, &listOut, &sError)

			if len(listOut) < 1 {
				sError = "報價失敗，查詢不到該訂單"
				bOk = false
			} else {
				C.InterfaceToStruct(listOut[0], &current)
				bOk = true

			}

		}

	}

	if bOk && order.Step != "0" {

		if order.Step == "1" {
			bOk, sError = orderStep1(current, &order)
		} else if order.Step == "2" {
			bOk, sError = orderStep2(current, &order)
		} else if order.Step == "3" {
			bOk, sError = orderStep3(current, &order)
		} else if order.Step == "4" {
			bOk, sError = orderStep4(current, &order)
		} else if order.Step == "5" {
			bOk = true
		} else if order.Step == "-1" {
			bOk, sError = orderCancel(current, &order)
		}

	}

	if bOk {
		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.OrderData(), C.StructToMap(order), &sError, true)
		sOkMsg = "訂單送出"
	}

	return
}

func UpdateQueryCount(Data C.VariantMap) {
	tmpMap := make(map[string]interface{})
	var tmpList = []interface{}{}
	var sError string
	tmpMap["GameItemSid"] = Data["GameItemSid"]

	itemCount := C.DataItemCount{}
	C.MapToStruct(Data, &itemCount)

	dataQuery := C.DataQueryCount{
		GameSid:      itemCount.GameSid,
		GameItemSid:  itemCount.GameItemSid,
		Name:         itemCount.Name,
		TotalCount:   int64(itemCount.TotalCount),
		TotalSell:    int64(itemCount.TotalSell),
		CurrentCount: int64(itemCount.TotalCount - itemCount.TotalSell),
	}

	tmpOk := CSQL.QueryTb(C.SQL_TABLE.QueryCount(), tmpMap, &tmpList, &sError)
	if tmpOk {

		if len(tmpList) > 0 {

			CSQL.UpdateTb(C.SQL_TABLE.QueryCount(), tmpMap, C.StructToMap(dataQuery), &sError)
		} else {

			CSQL.InsertTb(C.SQL_TABLE.QueryCount(), C.StructToMap(dataQuery), &sError, true)
		}

	}
}

func updateItemPrice(sGameSid string, gameRate string) {

	fRate, _ := strconv.ParseFloat(gameRate, 64)

	in := make(map[string]interface{})
	in["GameSid"] = sGameSid

	listOut := []interface{}{}
	var sError string
	CSQL.QueryTb(C.SQL_TABLE.GameItem(), in, &listOut, &sError)

	for _, v := range listOut {
		var item C.DataGameItem
		C.InterfaceToStruct(v, &item)

		f, err := strconv.ParseFloat(item.Bonus, 64)
		if err != nil {
			fmt.Println("Error parsing string:", err)
			continue
		} else {
			rounded := math.Ceil(f * fRate)
			s := strconv.FormatFloat(rounded, 'f', -1, 64)
			item.NTD = s
			//do something
		}

	}

}

func getCustomer(sSid string, data *C.CustomerData) bool {
	in := make(map[string]interface{})
	var listOut []interface{}
	in["Sid"] = sSid
	var sError string
	bOk := CSQL.QueryTb(C.SQL_TABLE.CustomerData(), in, &listOut, &sError)
	if !bOk || len(listOut) < 1 {
		return false
	}

	C.InterfaceToStruct(listOut[0], data)
	return true
}

func checkItemCount(orderData *C.OrderData, sErrorGameItemSid []string) ([]C.DataItemCount, bool) {

	var listLast []C.DataItemCount

	fmt.Println("checkItemCount : ")

	if orderData == nil {
		fmt.Println("order data is nil")
		return listLast, false
	}

	listItem := orderData.GetList("Item")

	for _, item := range listItem {
		parts := strings.Split(item, ",,")
		gameItemSid := parts[0]
		iCount, _ := strconv.ParseInt(parts[1], 10, 64)

		in := make(map[string]interface{})
		var out []interface{}
		in["GameItemSid"] = gameItemSid
		in["DESC"] = "Sid"
		in["LIMIT"] = "1"
		var sError string
		CSQL.QueryTb(C.SQL_TABLE.GameItemCount(), in, &out, &sError)

		itemCount := C.DataItemCount{
			GameItemSid: gameItemSid,
			ChangeValue: 0,
			TotalCount:  0,
			TotalSell:   0,
		}

		if len(out) > 0 {

			err := C.InterfaceToStruct(out[0], &itemCount)
			if err != nil {
				sError = err.Error()
				fmt.Printf("err : %s \n", err.Error())
				return listLast, false
			}
		}

		listLast = append(listLast, itemCount)

		iNowCount := itemCount.TotalCount - itemCount.TotalSell
		if iCount > int64(iNowCount) {
			sErrorGameItemSid = append(sErrorGameItemSid, gameItemSid)
			return listLast, false
		}
	}

	return listLast, true
}

func isBackSayCost(orderData C.OrderData) bool {
	// 排除返回未处理
	tmp0 := make(map[string]interface{})
	var listTmp []interface{}
	tmp0["Sid"] = orderData.Sid
	tmp0["Step"] = "2"
	var sError string
	bOk := CSQL.QueryTb(C.SQL_TABLE.OrderData(), tmp0, &listTmp, &sError) // SQL_TABLE_OrderData是SQL_TABLE::OrderData()的Go表示
	if !bOk {
		return false
	}

	return len(listTmp) > 0
}

func orderStep0(order *C.OrderData) (bool, string) {
	var cus C.CustomerData

	if !getCustomer(order.CustomerSid, &cus) {
		return false, "報價失敗，查詢客戶資料錯誤。"
	}

	var listSt []string

	_, bOk := checkItemCount(order, listSt)
	if !bOk {
		return false, "報價失敗, 商品庫存數量不足。"
	}

	order.Currency = cus.Currency
	order.CustomerName = cus.Name

	in := map[string]interface{}{
		"Sid": order.GameSid,
	}
	var tmpOut []interface{}
	var sError string

	if !CSQL.QueryTb(C.SQL_TABLE.GameList(), in, &tmpOut, &sError) {
		return false, sError
	}

	if len(tmpOut) > 0 {

		var tmpMap = make(map[string]interface{})
		C.InterFaceToMap(tmpOut[0], &tmpMap)

		order.GameRate = tmpMap["GameRate"].(string)
	} else {
		return false, "報價失敗，查詢遊戲資料錯誤"
	}
	listMoney := order.GetList("Money")
	if strings.ToUpper(order.Currency) == "NTD" {
		//	order.Money[0] = order.Cost
		listMoney[0] = "NTD"
	}

	cost, err := strconv.ParseFloat(order.Cost, 64) // 假设这个转换是安全的
	if err == nil && cost == 0 {
		//	order.Money[0] = "0"
		listMoney[0] = "0"
	}

	order.SetList("Money", listMoney)

	if len(order.StepTime) > 0 {

		list := order.GetList("StepTime")
	
		list[0] = C.TimeUtc8Str()

		order.SetList("StepTime", list)
	}

	return true, ""
}

func orderCancel(current C.OrderData, order *C.OrderData) (bool, string) {
	fmt.Println("CCCCC : order cancal")
	var sError string
	if current.Step != "0" {
		changeItemCount(current, true, &sError)
	}

	return true, ""
}

func orderStep1(current C.OrderData, order *C.OrderData) (bool, string) {
	bOnlyChangeOwner := false
	isBackSayCost := false
	if current.Owner == "未分配" || strings.TrimSpace(current.Owner) == "" {
		bOnlyChangeOwner = true
	} else if current.Step == "2" {
		isBackSayCost = true
	}

	var sError string
	if !isBackSayCost {
		if !changeItemCount(*order, false, &sError) {
			return false, "下單失敗, 商品庫存數量不足。"
		} else if current.Note1 == order.Note1 && current.Step != "0" && !bOnlyChangeOwner {
			return false, "下單失敗，目標訂單不處於報價狀態。"
		} else {

			if len(order.Id) < 6 {
				order.Id, _ = getNewOrderId(order.OrderDate)
			}

			listStepTime := C.StringToList(order.StepTime, ",", 6)
			listStepTime[1] = C.TimeUtc8Str()
			if strings.TrimSpace(listStepTime[0]) == "" {
				listStepTime[0] = C.TimeUtc8Str()
			}

			order.StepTime = C.ListToString(listStepTime, ",")

			in := map[string]interface{}{
				"Owner":     order.Owner,
				"OrderDate": order.OrderDate,
			}
			var tmpOut []interface{}
			var sError string

			CSQL.QueryTb(C.SQL_TABLE.OrderData(), in, &tmpOut, &sError)
			iSeq := len(tmpOut)
			sDash := ""
			if !strings.HasSuffix(order.Owner, "-") {
				sDash = "-"
			}

			if current.Sid == order.Sid && current.Step == "1" && !bOnlyChangeOwner {
				// 改留言，不重新給與編號
			} else {
				if order.Owner != "未分配" {
					// iSeq should be defined somewhere above.
					order.Name = order.Owner + sDash + fmt.Sprintf("%02d", iSeq+1)
				}
			}

		}
	}

	return true, ""
}

func orderStep2(current C.OrderData, order *C.OrderData) (bool, string) {

	listUser := C.StringToList(current.User, ",", 6)
	if current.Step == "2" {
		if len(order.PaddingUser) > 0 && string(listUser[2]) != order.PaddingUser {
			return false, GetUser(string(listUser[2])).Name + " 已接單處理"
		}
	}

	if current.Step == "3" {

		return false, GetUser(string(listUser[3])).Name + " 已儲值完成"
	}

	return true, ""
}

func orderStep3(current C.OrderData, order *C.OrderData) (bool, string) {

	listUser := C.StringToList(current.User, ",", 6)

	if order.Step == "3" && len(order.PaddingUser) > 0 {

		if current.Step == "3" && len(current.PaddingUser) > 0 {
			return false, GetUser(string(current.PaddingUser)).Name + " 正在回報中"
		}

		if current.Step == "4" {

			return false, GetUser(string(listUser[4])).Name + " 已回報"
		}
	}

	return true, ""
}

func orderStep4(current C.OrderData, order *C.OrderData) (bool, string) {

	var cus C.CustomerData = GetCustomer(order.CustomerSid)
	in := make(map[string]interface{})
	in["CustomerSid"] = order.CustomerSid
	in["DESC"] = "Sid"

	var preTotal float64
	listOut := []interface{}{}
	var sError string
	CSQL.QueryTb(C.SQL_TABLE.CustomerCost(), in, &listOut, &sError)

	if len(listOut) > 0 {
		var preCost C.CustomerCost
		C.InterfaceToStruct(listOut[0], &preCost)

		value, err := strconv.ParseFloat(preCost.Total, 64) // 第二個參數64表示轉換成float64
		if err == nil {
			preTotal = value
		}
	}

	var cost C.CustomerCost
	cost.UserSid = C.StringToList(order.User, ",", 6)[3]
	cost.Rate = order.ExRateSid
	cost.CustomerSid = order.CustomerSid
	cost.IsAddCost = false
	cost.Currency = cus.Currency
	cost.ChangeValue = "-" + order.Cost
	preTotal = preTotal + C.StringToFloat64(cost.ChangeValue)
	cost.Total = C.Float64ToString(preTotal)
	cost.OrderId = order.Id
	cost.OrderTime = order.OrderDate + order.OrderTime
	CSQL.InsertTb(C.SQL_TABLE.CustomerCost(), C.StructToMap(cost), &sError, false)
	changeMoney(cus, cost.Total, &sError)

	return true, ""
}

func changeItemCount(orderData C.OrderData, bIsAdd bool, sError *string) bool {
	fmt.Println("changeItemCount")
	*sError = ""
	var listLast []C.DataItemCount
	var listSt []string
	var bOk bool
	listLast, bOk = checkItemCount(&orderData, listSt)
	if !bIsAdd && !bOk {
		*sError = "商品庫存數量不足。"
		fmt.Println("商品庫存數量不足")
		return false

	}

	s := orderData.Item

	// 切割字串為二維字串陣列
	pairs := strings.Split(s, ";;")

	for i, pair := range pairs {
		v := strings.Split(pair, ",,")
		if i >= len(listLast) {
			continue
		}

		itemSid := v[0]
		count := v[1]

		var item C.DataGameItem = GetGameItem(itemSid)

		var itemCount C.DataItemCount
		itemCount.OrderSid = strconv.FormatInt(int64(orderData.Sid), 10)
		itemCount.GameSid = item.GameSid
		itemCount.GameItemSid = itemSid
		itemCount.Name = item.Name

		itemCount.ChangeValue = C.StringToInt64(count)

		if !bIsAdd {
			itemCount.ChangeValue *= -1
		}

		users := orderData.User

		var listUser []string = C.StringToList(users, ",", 6)

		itemCount.UserSid = listUser[1]

		itemCount.TotalCount = listLast[i].TotalCount

		itemCount.TotalSell = listLast[i].TotalSell + (itemCount.ChangeValue * -1)

		bOk, _, dataCount := CSQL.InsertTb(C.SQL_TABLE.GameItemCount(), C.StructToMap(itemCount), sError, true)

		if value, ok := dataCount["TotalCount"].(float64); ok {
			dataCount["TotalCount"] = int64(value)
		} else {
			// Handle error or situation where "TotalCount" is not a float64
		}

		if value, ok := dataCount["TotalSell"].(float64); ok {
			dataCount["TotalSell"] = int64(value)
		} else {
			// Handle error or situation where "TotalCount" is not a float64
		}

		if value, ok := dataCount["ChangeValue"].(float64); ok {
			dataCount["ChangeValue"] = int64(value)
		} else {
			// Handle error or situation where "TotalCount" is not a float64
		}

		if bOk {
			UpdateQueryCount(dataCount)
		} else {
			return false
		}

	}

	return true
}

func getNewOrderId(sOrderDate string) (string, error) {
	orderIdAdd := func(sDate string, last string) string {
		date := sDate

		if len(date) > 6 {
			date = date[len(date)-6:]
		}

		if last == "" {
			return date + "-" + "A001"
		}

		var sId string
		if len(last) >= 4 {
			sTmp := last[len(last)-4:]

			sSecond := sTmp[1:4]
			sFirst := sTmp[0:1]

			secondInt, err := strconv.Atoi(sSecond)
			if err != nil {
				// handle error
				return ""
			}

			if secondInt < 999 {
				sNum := fmt.Sprintf("%03d", secondInt+1)
				sId = date + "-" + sFirst + sNum
			} else {
				nextChar := rune(sFirst[0]) + 1
				sId = date + "-" + string(nextChar) + "001"
			}
		}

		return sId
	}

	sTodayLast, err := CSQL.LastOrderId(sOrderDate)
	if err != nil {

		return "", err
	}

	sRe := orderIdAdd(sOrderDate, sTodayLast)

	in := map[string]interface{}{
		"Id": sRe,
	}
	var tmpOut []interface{}
	var sError string

	if CSQL.QueryTb(C.SQL_TABLE.OrderData(), in, &tmpOut, &sError) && len(tmpOut) > 0 {
		sRe = orderIdAdd(sOrderDate, sRe)
	}

	return sRe, nil
}

func changeMoney(cus C.CustomerData, sValue string, sError *string) bool {
	in := map[string]interface{}{
		"Sid": cus.Sid,
	}

	var listOut []interface{}

	b := CSQL.QueryTb(C.SQL_TABLE.CustomerMoney(), in, &listOut, sError)

	if !b {
		return false
	}

	var money C.CustomerMoney

	if len(listOut) > 0 {
		C.InterfaceToStruct(listOut[0], &money)
	} else {
		money.Sid = cus.Sid
		money.Id = cus.Id
		money.Name = cus.Name
		money.Currency = cus.Currency
	}

	money.Money = sValue

	bOk, _, _ := CSQL.InsertTb(C.SQL_TABLE.CustomerMoney(), C.StructToMap(money), sError, true)

	if !bOk {
		return false
	}

	return true
}
