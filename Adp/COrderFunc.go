package Adp

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

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
	C.MapToStruct(Data, &order)
	if order.Step == "0" {
		bOk, sError = orderStep0(&order)
		sOkMsg = "報價成功"
	} else {
		if order.Owner == "未分配" {
			current = order
		} else {
			tmpIn["Sid"] = order.Sid
			var listOut []interface{}
			CSQL.QueryTb(C.SQL_TABLE.OrderData(), tmpIn, &listOut, &sError)

			if len(listOut) < 1 {
				sError = "報價失敗，查詢不到該訂單"
				bOk = false
			} else {
				C.InterfaceToStruct(listOut[0], &current)
				bOk = true

				if order.Step == "-1" {
					//to do
				} else if order.Step == "1" {
					//bOk,sError = orderStep1(current,&order)
				}

			}

		}

	}

	bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.OrderData(), C.StructToMap(order), &sError, true)
	sOkMsg = "訂單送出"

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
		TotalCount:   itemCount.TotalCount,
		TotalSell:    itemCount.TotalSell,
		CurrentCount: itemCount.TotalCount - itemCount.TotalSell,
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

func checkItemCount(orderData *C.OrderData, listLast []C.DataItemCount, sErrorGameItemSid []string) bool {

	if orderData == nil {
		fmt.Println("order data is nil")
		return false
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
				return false
			}
		}

		listLast = append(listLast, itemCount)
		iNowCount := itemCount.TotalCount - itemCount.TotalSell
		if iCount > iNowCount {
			sErrorGameItemSid = append(sErrorGameItemSid, gameItemSid)
			return false
		}
	}

	return true
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

	var listLast []C.DataItemCount
	var listSt []string
	if !checkItemCount(order, listLast, listSt) {
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
		list[0] = time.Now().UTC().Add(time.Minute * 8).Format("20060102150405")
		order.SetList("StepTime", list)
	}

	return true, ""
}

/*
func orderStep1(current C.OrderData,order *C.OrderData) (bool, string) {
	bOnlyChangeOwner := false

	if current.Owner == "未分配" || strings.TrimSpace(current.Owner) == "" {
		bOnlyChangeOwner = true
	}

	if !isBackSayCost(*order) {
		if !changeItemCount(order, false) {
			return false, "下單失敗, 商品庫存數量不足。"
		} else if current.Note1 == order.Note1 && current.Step != "0" && !bOnlyChangeOwner {
			return false, "下單失敗，目標訂單不處於報價狀態。"
		} else {
			checkUpdate(ACT_ADD_ITEM_COUNT)

			if strings.TrimSpace(order.Id) == "" {
				order.Id = getNewOrderId(order.OrderDate)
			}

			if len(order.StepTime) > 1 {
				now := time.Now().Add(8 * time.Hour)
				order.StepTime[1] = now.Format("20060102150405")
				if strings.TrimSpace(order.StepTime[0]) == "" {
					order.StepTime[0] = now.Format("20060102150405")
				}
			}

			// 這裡是Qt特定的代碼，您需要找到Go中的替代方法
			// ...

			if strings.TrimSpace(order.Owner) == "" {
				order.Owner = "None"
			}
			// iSeq := ...  // Again, this is a Qt method. You need a Go equivalent.

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
*/
