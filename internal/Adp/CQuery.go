package Adp

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strconv"

	C "adpGo/common"
	CSQL "adpGo/pkg/CSql"
	NETWORK "adpGo/pkg/RpkNetwork"
)

//show global status like 'com_stmt%';

var LoginUser = make(map[string]interface{})

type AdpRecaller struct{}

func (d AdpRecaller) ImplementRecall(data NETWORK.CData) NETWORK.CData {

	var reData C.VariantMap = make(map[string]interface{})
	reList := []interface{}{}
	iAction := data.Action
	var sError, sOkMsg string
	var bOk bool = false

	var tmpIn C.VariantMap = make(map[string]interface{})

	var tmpMap C.VariantMap = make(map[string]interface{})
	// var tmpList = []interface{}{}
	var Data C.VariantMap = data.Data

	switch {
	case iAction >= C.API_REQUSET && iAction < 9900:
		fmt.Println("api")

	case iAction == 1:
		/*
			var d2 interface{} = data.Data
			var d3 map[string]interface{} = d2.(map[string]interface{})
			userSid, _ := d3["UserSid"].(string)
			Session, _ := d3["Session"].(string)
			bOk = LoginUser[userSid] == Session
		*/
		//to do

		bOk = true

	case iAction == C.QUERY_INFO:
		tmpMap["ServerVersion"] = "1.5"
		reData = tmpMap
		bOk = true

	case iAction == C.LOGIN:

		C.PrintHHMMSS("A0")

		fmt.Println("login")
		//conditions:=[string]interface
		//CSQL.QueryTb(tableName, conditions, listOut, sError)
		bOk = false
		sError = "連線失敗"

		bCheckVer := false

		if len(data.ListData) >= 3 {

			in := map[string]interface{}{
				"Id":       data.ListData[0],
				"Password": data.ListData[1],
			}

			verStr, _ := data.ListData[2].(string)
			fmt.Println("解析版本字串：", verStr)
			ver, err := strconv.ParseFloat(verStr, 64)
			if err != nil {
				fmt.Println("解析版本號錯誤：", err)

			} else {
				if ver > 2.3 {
					fmt.Println("版本大於 2.3")
					bCheckVer = true
				} else {
					sError = "版本不符合"
				}
			}

			if bCheckVer {

				tmpB := CSQL.QueryTb(C.SQL_TABLE.UserData(), in, &reList, &sError)

				if tmpB {
					sError = "帳密錯誤"
				}

				if len(reList) > 0 {

					strBytes := []byte(C.TimeUtc8Str())
					hash := md5.Sum(strBytes)
					sSession := hex.EncodeToString(hash[:])

					reList[0].(map[string]interface{})["Session"] = sSession
					fmt.Println(sSession)
					C.InterFaceToMap(reList[0], reData.Origin())

					var UserSid int
					if sid, ok := reData["Sid"].(float64); ok {
						UserSid = int(sid)
					} else {
						// 处理无法转换为字符串的情况
						fmt.Println("Failed to convert Sid to string")
					}

					LoginUser[strconv.Itoa(UserSid)] = sSession
					bOk = true
					sOkMsg = "登入成功"

					fmt.Println("loging msg : " + sOkMsg)
				}

			}
		}

		//
	case iAction == C.SET_VALUE:

		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.Settings(), Data, &sError, true)

	//QUERY============

	case iAction == C.GET_VALUE:

		bOk = CSQL.QueryTb(C.SQL_TABLE.Settings(), Data, &reList, &sError)
		if bOk && len(reList) > 0 {
			C.InterFaceToMap(reList[0], reData.Origin())
		}

	case iAction == C.QUERY_USER:

		bOk = CSQL.QueryTb(C.SQL_TABLE.UserData(), Data, &reList, &sError)

	case iAction == C.QUERY_CUSTOMER:
		bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerData(), Data, &reList, &sError)

	case iAction == C.QUERY_CUSTOMER_GAME_INFO:
		bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerGameInfo(), Data, &reList, &sError)

	case iAction == C.QUERY_CUSTOMER_COST:
		bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerCost(), Data, &reList, &sError)

	case iAction == C.QUERY_GAME_LIST:
		bOk = CSQL.QueryTb(C.SQL_TABLE.GameList(), Data, &reList, &sError)

	case iAction == C.QUERY_GAME_RATE:

		bOk = CSQL.QueryTb(C.SQL_TABLE.GameRate(), Data, &reList, &sError)

		if bOk && len(reList) < 1 {
			tmpOut := make([]interface{}, 0)
			tmp := map[string]interface{}{
				"sid": Data["GameSid"],
			}

			bOk = CSQL.QueryTb(C.SQL_TABLE.GameList(), tmp, &tmpOut, &sError)

			if bOk && len(tmpOut) > 0 {

				var GameData C.DataGameList

				C.InterfaceToStruct(tmpOut[0], &GameData)

				var gameRateData C.DataGameRate
				gameRateData.GameSid = strconv.Itoa(GameData.Sid)
				gameRateData.GameName = GameData.Name
				gameRateData.UserSid = GameData.UserSid
				gameRateData.Rate = GameData.GameRate

				var sErrorTmp string
				go CSQL.InsertTb(C.SQL_TABLE.GameRate(), C.StructToMap(gameRateData), &sErrorTmp, true)
			}

		}

	case iAction == C.QUERY_GAME_ITEM:
		bOk = CSQL.QueryTb(C.SQL_TABLE.GameItem(), Data, &reList, &sError)

	case iAction == C.QUERY_BULLETIN:
		bOk = CSQL.QueryTb(C.SQL_TABLE.Bulletin(), Data, &reList, &sError)

	case iAction == C.QUERY_CUSTOM_CLASS:
		in := make(map[string]interface{})
		in["Type"] = "group"
		bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerClass(), in, &reList, &sError)

	case iAction == C.QUERY_CUSTOM_DEBIT:
		in := make(map[string]interface{})
		in["Type"] = "debit"
		bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerClass(), in, &reList, &sError)

	case iAction == C.QUERY_FACTORY_CLASS:

		var d map[string]interface{}

		if len(data.ListData) > 0 {
			if tmp, ok := data.ListData[0].(map[string]interface{}); ok {
				d = tmp
			}
		}

		if d == nil {
			d = Data
		}

		bOk = CSQL.QueryTb(C.SQL_TABLE.FactoryClass(), d, &reList, &sError)
		/*
			fac := C.DataFactory{
				Sid:  999,
				Id:   "未",
				Name: "未分配",
			}
			reList = append(reList, fac)
		*/

	case iAction == C.QUERY_PAY_TYPE:
		bOk = CSQL.QueryTb(C.SQL_TABLE.PayType(), Data, &reList, &sError)

	// case iAction == C.QUERY_GROUP:
	// 	bOk = CSQL.QueryTb(C.SQL_TABLE.GroupData(), Data, &reList, &sError)

	case iAction == C.QUERY_ORDER:
		bOk = CSQL.QueryTb(C.SQL_TABLE.OrderData(), Data, &reList, &sError)

	case iAction == C.QUERY_BOUNS:
		bOk = CSQL.QueryTb(C.SQL_TABLE.UserBonus(), Data, &reList, &sError)

	case iAction == C.QUERY_SCHEDULE:
		bOk = CSQL.QueryTb(C.SQL_TABLE.Schedule(), Data, &reList, &sError)

	case iAction == C.QUERY_EXCHANGE:
		bOk = CSQL.QueryTb(C.SQL_TABLE.ExchangeRate(), Data, &reList, &sError)

	case iAction == C.QUERY_PRIMERATE:
		bOk = CSQL.QueryTb(C.SQL_TABLE.PrimeCostRate(), Data, &reList, &sError)

	case iAction == C.QUERY_PIC:
		bOk = CSQL.QueryTb(C.SQL_TABLE.PicData(), Data, &reList, &sError)

	case iAction == C.QUERY_ITEM_COUNT:
		bOk = CSQL.QueryTb(C.SQL_TABLE.GameItemCount(), Data, &reList, &sError)

	case iAction == C.QUERY_DEBIT_CLASS:
		tmpIn := make(map[string]interface{})
		tmpIn["ASC"] = "Sort"

		bOk = CSQL.QueryTb(C.SQL_TABLE.DebitClass(), tmpIn, &reList, &sError)
	case iAction == C.QUERY_COUNT:

		bOk = CSQL.QueryTb(C.SQL_TABLE.QueryCount(), Data, &reList, &sError)
	case iAction == C.QUERY_MIX:

		for key, value := range Data {
			//fmt.Printf("key: %s \n ", key)
			api, err := strconv.Atoi(key)

			if err != nil || api == C.QUERY_MIX {
				continue
			}
			//fmt.Printf("api: %d \n ", api)

			buff := []byte(value.(string))
			//fmt.Printf("buffLen: %d \n", len(buff))

			var tmp NETWORK.CData

			if err := tmp.DecodeJSON(buff); err != nil {
				continue
			}

			tmp.State = NETWORK.ACT_LOCAL
			tmpRe := d.ImplementRecall(tmp)

			if tmp2, err := tmpRe.EncodeJSON(); err == nil {
				//fmt.Printf("%d encode ok ", api)
				reData[key] = tmp2
			}

		}

	case iAction == C.QUERY_DAY_REPORT:
	//reList =  GetDayReport(Data)

	//END QUERY--------------------------------

	case iAction == C.ADD_USER:

		if len(data.ListData) < 6 {
			break
		}

		user := C.UserData{
			Id:       C.InterfaceToString(data.ListData[0]),
			Password: data.ListData[1].(string),
			Cid:      data.ListData[2].(string),
			Name:     data.ListData[3].(string),
			// 對於Lv，因為它是int型別，所以您需要進行適當的轉換。以下只是一個示範：
			Lv:       C.InterfaceToInt(data.ListData[4]),
			StartDay: data.ListData[5].(string),
			// 假設CreateTime是由其他函數生成，所以這裡不直接從ListData取值
			CreateTime: CSQL.CurrentTime(),
			// 以下是其他可能的欄位賦值
		}

		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.UserData(), C.StructToMap(user), &sError, false)
		sOkMsg = "新增成功"

	case iAction == C.EDIT_USER:
		in := make(map[string]interface{})

		_, ok := Data["Sid"]
		if !ok {
			sError = "資料異常：缺少 Sid"
			break
		}

		in["Sid"] = Data["Sid"]
		bOk = CSQL.UpdateTb(C.SQL_TABLE.UserData(), in, Data, &sError)
		sOkMsg = "修改成功"

	case iAction == C.DEL_USER:
		bOk = CSQL.DelFromTb(C.SQL_TABLE.UserData(), Data, &sError)
		sOkMsg = "刪除成功"

	case iAction == C.ADD_CUSTOMER:

		bOk, _, tmpMap = CSQL.InsertTb(C.SQL_TABLE.CustomerData(), Data, &sError, false)

		sOkMsg = "客戶名稱:" + tmpMap.String("Name") + "\n新增成功"

	case iAction == C.EDIT_CUSTOMER:

		// 確保 Sid 存在
		_, ok := Data["Sid"]
		if !ok {
			sError = "資料異常：缺少 Sid"
			break
		}

		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.UpdateTb(C.SQL_TABLE.CustomerData(), tmpIn, Data, &sError)
		sOkMsg = "修改成功"

	case iAction == C.DEL_CUSTOMER:

		_, ok := Data["Sid"]
		if !ok {
			sError = "資料異常：缺少 Sid"
			break
		}
		in := make(map[string]interface{})
		in["Sid"] = Data["Sid"]
		bOk = CSQL.DelFromTb(C.SQL_TABLE.CustomerData(), in, &sError)
		if bOk {
			CSQL.DelFromTb(C.SQL_TABLE.CustomerMoney(), in, &sError)

			in = make(map[string]interface{})
			in["CustomerSid"] = Data["Sid"]
			CSQL.DelFromTb(C.SQL_TABLE.CustomerGameInfo(), in, &sError)

			CSQL.DelFromTb(C.SQL_TABLE.CustomerCost(), in, &sError)
		}

		sOkMsg = "刪除成功"

	case iAction == C.REPLACE_GAME_INFO:
		bOk = true
		for i := 0; i < len(data.ListData); i++ {
			bTmp, _, _ := CSQL.InsertTb(C.SQL_TABLE.CustomerGameInfo(), data.ListData[i].(map[string]interface{}), &sError, true)

			if !bTmp {
				bOk = false
			}
		}

		sOkMsg = "客戶遊戲資料修改完成"

	case iAction == C.DEL_GAME_INFO:
		bOk = true
		for i := 0; i < len(data.ListData); i++ {
			C.InterFaceToMap(data.ListData[i], tmpMap.Origin())
			tmpIn["Sid"] = tmpMap.String("Sid")
			bTmp := CSQL.DelFromTb(C.SQL_TABLE.CustomerGameInfo(), tmpIn, &sError)

			if !bTmp {
				bOk = false
			}
		}

		sOkMsg = "客戶遊戲資料修改完成"

	case iAction == C.ADD_CUSTOMER_COST:
		bOk, _, tmpMap = CSQL.InsertTb(C.SQL_TABLE.CustomerCost(), Data, &sError, false)

		sOkMsg = "加值完成"
	case iAction == C.LAST_CUSTOMER_COST_ID:
		var tmpId string

		tmpId, bOk = CSQL.LastCustomerAddCostID()
		tmpMap["OrderId"] = tmpId
		reList = append(reList, tmpMap)

	case iAction == C.ADD_GAME_LIST:
		var Sid int64
		bOk, Sid, _ = CSQL.InsertTb(C.SQL_TABLE.GameList(), Data, &sError, false)

		gameList := C.DataGameList{}
		C.MapToStruct(Data, &gameList)
		gameRate := C.DataGameRate{}
		gameRate.GameSid = strconv.FormatInt(Sid, 10)
		gameRate.Rate = gameList.GameRate
		gameRate.GameName = gameList.Name
		gameRate.UserSid = gameList.UserSid

		tmpMap = C.StructToMap(gameRate)
		if bOk {
			CSQL.InsertTb(C.SQL_TABLE.GameRate(), tmpMap, &sError, false)
		}
		sOkMsg = "遊戲新增成功"

	case iAction == C.EDIT_GAME_LIST:

		_, ok := Data["Sid"]
		if !ok {
			sError = "資料異常：缺少 Sid"
			break
		}

		in := map[string]interface{}{
			"Sid": Data["Sid"],
		}

		var game C.DataGameList

		C.SetData(Data, &game)

		//tmpRate := strconv.FormatFloat(game.GameRate, 'f', 2, 64) // 'f' 表示格式，2 表示小數點後的位數，64 表示它是 float64

		gameRate := map[string]interface{}{
			"GameSid":  game.Sid,
			"GameName": game.Name,
			"UserSid":  game.UserSid,
			"Rate":     game.GameRate,
		}

		bOk = CSQL.UpdateTb(C.SQL_TABLE.GameList(), in, C.GetData(&game), &sError)

		if bOk {
			bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.GameRate(), gameRate, &sError, true)
			//game item NTD 沒有使用，客戶端直接用GameRate*Bonus計算出來顯示
			//updateItemPrice(game.Sid,game.GameRate)
		}
		sOkMsg = "遊戲修改完成"

	case iAction == C.DEL_GAME_LIST:

		_, ok := Data["Sid"]
		if !ok {
			sError = "資料異常：缺少 Sid"
			break
		}

		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.DelFromTb(C.SQL_TABLE.GameList(), tmpIn, &sError)
		if bOk {
			tmpMap["GameSid"] = tmpIn["Sid"]
			CSQL.DelFromTb(C.SQL_TABLE.GameRate(), tmpMap, &sError)
			CSQL.DelFromTb(C.SQL_TABLE.GameItem(), tmpMap, &sError)
			CSQL.DelFromTb(C.SQL_TABLE.GameItemCount(), tmpMap, &sError)
			CSQL.DelFromTb(C.SQL_TABLE.QueryCount(), tmpMap, &sError)

		}
		sOkMsg = "遊戲刪除成功"

	case iAction == C.ADD_GAME_ITEM:

		var iCount int64 = 0

		var bHasCount bool = false
		if count, ok := Data["Count"]; ok {

			switch v := count.(type) {
			case int:
				iCount = int64(v)
			case int32:
				iCount = int64(v)
			case int64:
				iCount = v
			case float64:
				iCount = int64(v)

			}

			bHasCount = true
			// 移除該 key
			delete(Data, "Count")
		}

		var itemSid int64
		bOk, itemSid, _ = CSQL.InsertTb(C.SQL_TABLE.GameItem(), Data, &sError, false)

		if bOk && bHasCount {
			item := C.DataGameItem{}
			C.MapToStruct(Data, &item)

			itemCount := C.DataItemCount{
				GameSid:     item.GameSid,
				GameItemSid: strconv.FormatInt(itemSid, 10),
				Name:        item.Name,
				TotalCount:  iCount,
			}

			tmpMap = C.StructToMap(itemCount)

			dataCount := make(map[string]interface{})
			bOk, _, dataCount = CSQL.InsertTb(C.SQL_TABLE.GameItemCount(), tmpMap, &sError, false)

			UpdateQueryCount(dataCount)
		}

		sOkMsg = "商品新增成功"

	case iAction == C.DEL_GAME_ITEM:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.DelFromTb(C.SQL_TABLE.GameItem(), tmpIn, &sError)
		if bOk {
			tmpMap["GameItemSid"] = Data["Sid"]
			CSQL.DelFromTb(C.SQL_TABLE.GameItemCount(), tmpMap, &sError)
			CSQL.DelFromTb(C.SQL_TABLE.QueryCount(), tmpMap, &sError)
		}

		sOkMsg = "商品刪除成功"

	case iAction == C.EDIT_GAME_ITEM:
		fmt.Println("Edit game item:")

		bHasList := len(data.ListData) > 0

		bOk = true

		if !bHasList {
			d := make(map[string]interface{})
			d["Sid"] = Data["Sid"]
			Data["Count"] = nil
			bOk = CSQL.UpdateTb(C.SQL_TABLE.GameItem(), d, Data, &sError)

		} else {

			bIsBatch := false

			if bIsBatch {

				//方式1
				var listConditions []map[string]interface{}

				for _, v := range data.ListData {
					var mapData map[string]interface{}
					C.InterFaceToMap(v, &mapData)

					if sidValue, ok := mapData["Sid"]; ok {
						condition := map[string]interface{}{
							"Sid": sidValue,
						}
						listConditions = append(listConditions, condition)
					}
				}

				listMap, _ := C.ToListMap(data.ListData)
				bOk = CSQL.BatchUpdateTb(C.SQL_TABLE.GameItem(), listConditions, listMap, &sError)

			} else {

				//方式2
				for _, t := range data.ListData {
					tmpMap, ok := t.(map[string]interface{}) // 這裡進行類型斷言
					if !ok {
						fmt.Println("t 不是 map[string]interface{} 的型別")
						continue // 如果斷言失敗，則跳過此次迴圈的迭代
					}
					tmp := make(map[string]interface{})
					tmp["Sid"] = tmpMap["Sid"]
					tmp["Count"] = nil
					// 假設 CSQL.UpdateTb 是一個用來更新數據表的Go函數
					b := CSQL.UpdateTb("GameItem", tmp, tmpMap, &sError)
					if !b {
						bOk = false
					}
				}

			}

		}

		sOkMsg = "商品修改完成"
		fmt.Println(sOkMsg)

	case iAction == C.ADD_ITEM_COUNT:

		tmp := make(map[string]interface{})
		tmp["GameItemSid"] = Data["GameItemSid"]
		tmp["DESC"] = "Sid"
		tmp["Limit"] = "1"
		var listOut []interface{}
		var tmpError string
		tmpOk := CSQL.QueryTb(C.SQL_TABLE.GameItemCount(), tmp, &listOut, &tmpError)

		if tmpOk && len(listOut) > 0 {
			var d C.DataItemCount
			C.InterfaceToStruct(listOut[0], &d)
			Data["TotalSell"] = d.TotalSell
			//該接口只用來 調整庫存
			//防止同步時間差導致賣出數量出錯
		}

		GameItemSid, _ := Data["GameItemSid"].(string)

		value, exists := Data["GameSid"]
		if !exists {
			// GameSid 在 map tmp 中不存在
			var item C.DataGameItem = GetGameItem(GameItemSid)
			Data["GameSid"] = item.GameSid
			Data["Name"] = item.Name

		} else if strValue, isString := value.(string); isString && len(strValue) < 1 {
			// GameSid 存在，但其值的長度小於 1
			var item C.DataGameItem = GetGameItem(GameItemSid)
			Data["GameSid"] = item.GameSid
			Data["Name"] = item.Name

		}

		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.GameItemCount(), Data, &sError, false)
		sOkMsg = "新增成功"
		if bOk {
			UpdateQueryCount(Data)
		}

	case iAction == C.DEL_ITEM_COUNT:
		bOk = CSQL.DelFromTb(C.SQL_TABLE.GameItemCount(), Data, &sError)
		bOk = CSQL.DelFromTb(C.SQL_TABLE.QueryCount(), Data, &sError)
		sOkMsg = "刪除成功"

	case iAction == C.EDIT_ITEM_COUNT:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.UpdateTb(C.SQL_TABLE.GameItemCount(), tmpIn, Data, &sError)
		sOkMsg = "修改成功"
		if bOk {
			go UpdateQueryCount(Data)
		}

	case iAction == C.ADD_BULLETIN:
		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.Bulletin(), Data, &sError, false)
		sOkMsg = "新增成功"

	case iAction == C.DEL_BULLETIN:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.DelFromTb(C.SQL_TABLE.Bulletin(), tmpIn, &sError)
		sOkMsg = "刪除成功"

	case iAction == C.EDIT_BULLETIN:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.UpdateTb(C.SQL_TABLE.Bulletin(), tmpIn, Data, &sError)
		sOkMsg = "修改完成"

	case iAction == C.ADD_CUSTOM_CLASS:
		Data["Type"] = "group"
		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.CustomerClass(), Data, &sError, false)
		sOkMsg = "新增成功"

	case iAction == C.DEL_CUSTOM_CLASS:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.DelFromTb(C.SQL_TABLE.CustomerClass(), tmpIn, &sError)
		sOkMsg = "刪除成功"

	case iAction == C.EDIT_CUSTOM_CLASS:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.UpdateTb(C.SQL_TABLE.CustomerClass(), tmpIn, Data, &sError)
		sOkMsg = "修改完成"

	case iAction == C.ADD_FACTORY_CLASS:
		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.FactoryClass(), Data, &sError, false)
		sOkMsg = "新增成功"

	case iAction == C.DEL_FACTORY_CLASS:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.DelFromTb(C.SQL_TABLE.FactoryClass(), tmpIn, &sError)
		sOkMsg = "刪除成功"

	case iAction == C.EDIT_FACTORY_CLASS:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.UpdateTb(C.SQL_TABLE.FactoryClass(), tmpIn, Data, &sError)
		sOkMsg = "修改完成"

	case iAction == C.ADD_PAY_TYPE:
		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.PayType(), Data, &sError, false)
		sOkMsg = "新增成功"

	case iAction == C.DEL_PAY_TYPE:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.DelFromTb(C.SQL_TABLE.PayType(), tmpIn, &sError)
		sOkMsg = "刪除成功"

	case iAction == C.EDIT_PAY_TYPE:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.UpdateTb(C.SQL_TABLE.PayType(), tmpIn, Data, &sError)
		sOkMsg = "修改完成"

		// case iAction == C.ADD_GROUP:
		// 	bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.GroupData(), Data, &sError, false)
		// 	sOkMsg = "新增成功"

		// case iAction == C.DEL_GROUP:
		// 	tmpIn["Sid"] = Data["Sid"]
		// 	bOk = CSQL.DelFromTb(C.SQL_TABLE.GroupData(), tmpIn, &sError)
		// 	sOkMsg = "刪除成功"

		// case iAction == C.EDIT_GROUP:
		// 	tmpIn["Sid"] = Data["Sid"]
		// 	bOk = CSQL.UpdateTb(C.SQL_TABLE.GroupData(), tmpIn, Data, &sError)
		// 	sOkMsg = "修改完成"

	case iAction == C.ADD_BOUNS:
		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.UserBonus(), Data, &sError, false)
		sOkMsg = "修改完成"

	case iAction == C.ADD_SCHEDULE:
		tmpIn["Sid"] = Data["Sid"]
		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.Schedule(), Data, &sError, true)
		sOkMsg = "新增成功"

	case iAction == C.ADD_EXCHANGE:
		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.ExchangeRate(), Data, &sError, false)
		sOkMsg = "新增成功"

	case iAction == C.ADD_PRIMERATE:
		tmpIn["Sid"] = Data["Sid"]
		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.PrimeCostRate(), Data, &sError, false)
		sOkMsg = "新增成功"

	case iAction == C.UPLOAD_PIC:
		go UploadPic(Data)
		bOk = true

	case iAction == C.ADD_DEBIT_CLASS:
		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.DebitClass(), Data, &sError, false)
		sOkMsg = "支付管道新增完成"

	case iAction == C.DEL_DEBIT_CLASS:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.DelFromTb(C.SQL_TABLE.DebitClass(), tmpIn, &sError)
		sOkMsg = "支付管道刪除成功"

	case iAction == C.EDIT_DEBIT_CLASS:
		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.UpdateTb(C.SQL_TABLE.DebitClass(), tmpIn, Data, &sError)
		sOkMsg = "支付管道修改完成"

	case iAction == C.PAY_ADD_COST:

		var customer C.CustomerData
		C.InterfaceToStruct(Data["CustomerData"], &customer)

		var costData C.CustomerCost
		C.InterfaceToStruct(Data["CostData"], &costData)

		var iChangeValue float64 = 0
		var iOldTotal float64 = 0
		var iNewTotal float64 = 0
		var err error

		iChangeValue, err = strconv.ParseFloat(costData.ChangeValue, 64)
		if err != nil {
			break
		}

		tmpIn["CustomerSid"] = strconv.Itoa(customer.Sid)

		tmpIn["DESC"] = "Sid"
		tmpIn["Limit"] = "1"
		var listOut []interface{}
		CSQL.QueryTb(C.SQL_TABLE.CustomerCost(), tmpIn, &listOut, &sError)

		if len(listOut) > 0 {

			var tmp C.CustomerCost
			C.InterfaceToStruct(listOut[0], &tmp)
			iOldTotal, err = strconv.ParseFloat(tmp.Total, 64)
			if err != nil {
				break
			}

		}

		iNewTotal = iOldTotal + iChangeValue

		strValue := fmt.Sprintf("%.2f", iNewTotal)
		costData.Total = strValue

		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.CustomerCost(), C.StructToMap(costData), &sError, true)

		if bOk {
			var money C.CustomerMoney
			var tmp = make(map[string]interface{})
			tmp["Sid"] = customer.Sid
			var listTmp []interface{}
			CSQL.QueryTb(C.SQL_TABLE.CustomerMoney(), tmp, &listTmp, &sError)

			if len(listTmp) > 0 {
				C.InterfaceToStruct(listTmp[0], &money)
			} else {
				money.Sid = customer.Sid
				money.Name = customer.Id
				money.Currency = customer.Currency
			}

			money.Money = strValue

			CSQL.InsertTb(C.SQL_TABLE.CustomerMoney(), C.StructToMap(money), &sError, true)

		}

		// bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.DebitClass(), Data, &sError, false)

		sOkMsg = "加值完成\n\n加值金額:" + fmt.Sprintf("%.2f", iChangeValue) + "\n加值後金額:" + fmt.Sprintf("%.2f", iNewTotal)

	case iAction == C.LAST_ORDER_ID:

		sDate, ok := Data["OrderDate"].(string)
		if ok {
			sReId, err := CSQL.LastOrderId(sDate)
			if err != nil {
				sError = err.Error()
				break
			}

			tmp := make(map[string]interface{})
			tmp["Id"] = sReId

			reList = append(reList, tmp)
		}

	case iAction == C.LAST_ORDER_NAME:

		sDate, ok1 := Data["OrderDate"].(string)
		sOwner, ok2 := Data["Owner"].(string)

		if ok1 && ok2 {
			sReId, err := CSQL.LastOrderName(sOwner, sDate)
			if err == nil {
				tmp := make(map[string]interface{})
				tmp["Name"] = sReId
				reList = append(reList, tmp)
			} else {
				sError = err.Error()
				fmt.Println("Error:", err)
				break
			}
		}

	case iAction == C.UPDATE_DATA:

		sDate, ok1 := Data["OrderData"].(string)
		if ok1 {
			conditions := make(map[string]interface{})

			conditions["UpdateTime >"] = sDate
			conditions["LIMIT"] = "3000"
			conditions["OrderDate >="] = C.DateUtc8Str(-1)
			reOrder := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.OrderData(), conditions, &reOrder, &sError)
			reData["OrderData"] = reOrder
		}

		sDate, ok1 = Data["CustomerData"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["ASC"] = "UpdateTime"
			conditions["LIMIT"] = "5000"
			conditions["UpdateTime >"] = sDate
			reCus := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerData(), conditions, &reCus, &sError)

			condi := make(map[string]interface{})
			condi["UpdateTime ="] = sDate
			condi["ASC"] = "UpdateTime"
			reTmp := []interface{}{}
			CSQL.QueryTb(C.SQL_TABLE.CustomerData(), condi, &reTmp, &sError)
			if len(reTmp) > 1 {
				reTmp = reTmp[1:]               // 移除 reTmp 的第一個元素
				reCus = append(reCus, reTmp...) // 串連 reTmp 到 reCus
			}

			reData["CustomerData"] = reCus
		}

		sDate, ok1 = Data["UserData"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["UpdateTime >"] = sDate
			re := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.UserData(), conditions, &re, &sError)
			reData["UserData"] = re
		}

		sDate, ok1 = Data["GameList"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["UpdateTime >"] = sDate
			re := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.GameList(), conditions, &re, &sError)
			reData["GameList"] = re
		}

		sDate, ok1 = Data["GameItem"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["UpdateTime >"] = sDate
			re := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.GameItem(), conditions, &re, &sError)
			reData["GameItem"] = re
		}

		sDate, ok1 = Data["ExchangeRate"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["UpdateTime >"] = sDate
			re := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.ExchangeRate(), conditions, &re, &sError)
			reData["ExchangeRate"] = re
		}

		sDate, ok1 = Data["PrimeCostRate"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["UpdateTime >"] = sDate
			re := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.PrimeCostRate(), conditions, &re, &sError)
			reData["PrimeCostRate"] = re
		}

		sDate, ok1 = Data["CustomerClass"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["UpdateTime >"] = sDate
			re := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerClass(), conditions, &re, &sError)
			reData["CustomerClass"] = re
		}

		sDate, ok1 = Data["FactoryClass"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["UpdateTime >"] = sDate
			re := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.FactoryClass(), conditions, &re, &sError)

			/*
				fac := C.DataFactory{
					Sid:  999,
					Id:   "未",
					Name: "未分配",
				}
				re = append(re, fac)
			*/
			reData["FactoryClass"] = re
		}

		sDate, ok1 = Data["PayType"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["UpdateTime >"] = sDate
			re := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.PayType(), conditions, &re, &sError)
			reData["PayType"] = re
		}

		sDate, ok1 = Data["UserBonus"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["UpdateTime >"] = sDate
			re := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.UserBonus(), conditions, &re, &sError)
			reData["UserBonus"] = re
		}

		sDate, ok1 = Data["BulletinData"].(string)
		if ok1 {
			conditions := make(map[string]interface{})
			conditions["UpdateTime >"] = sDate
			conditions["DESC"] = "Sid"
			re := []interface{}{}
			bOk = CSQL.QueryTb(C.SQL_TABLE.Bulletin(), conditions, &re, &sError)
			reData["BulletinData"] = re
		}

	case iAction == C.REPLACE_ORDER, iAction == C.PAY_ORDER:

		bOk, sOkMsg, sError = DoOrder(data, &reData, &reList)

	case iAction == C.EDIT_SORT:
		sTable, ok := Data["Table"].(string)
		if ok && sTable == "FactoryClass" {

			// 確保 Table 的值為 "FactoryClass"，並且是 string 類型
			data, ok := Data["Data"].(map[string]interface{})
			if ok {
				// 確保 Data 的值是 map[string]interface{} 類型

				for key, value := range data {
					sSid := key // 因為 key 已經是 string 類型
					sSort, ok := value.(string)
					if ok {

						//fmt.Printf("AAAAAA:1 sSid: %s, sSort: %s\n", sSid, sSort)

						in := make(map[string]interface{})
						in["Sid"] = sSid
						conditions := make(map[string]interface{})
						conditions["Sort"] = sSort
						bOk = CSQL.UpdateTb(C.SQL_TABLE.FactoryClass(), in, conditions, &sError)

						if bOk {
							sOkMsg = "排序修改完成"
						}
					}
				}
			}
		}

	case iAction == C.RUN_SH:
		sRun, bCheck := Data["Run"].(string)

		if bCheck {
			cmd := exec.Command("/bin/sh", sRun)
			err := cmd.Run()
			if err != nil {
				bOk = false

			} else {
				bOk = true
			}

			if bOk {
				sOkMsg = "已啟動"
			}
		}

	default:
		// 未知操作，可以进行相应的处理
		fmt.Printf("unknown action : %d \n", iAction)

	}

	var re NETWORK.CData
	re.Action = data.Action
	re.State = NETWORK.ACT_RECALL

	re.Ok = bOk

	if re.Ok {
		re.Msg = sOkMsg
	} else {
		re.Msg = sError
	}

	re.Data = reData

	re.ListData = reList
	return re
}

func UploadPic(d map[string]interface{}) {
	var error string
	CSQL.InsertTb(C.SQL_TABLE.PicData(), d, &error, true)

}

func GetGameItem(sSid string) (re C.DataGameItem) {

	in := make(map[string]interface{})
	in["Sid"] = sSid

	listOut := []interface{}{}
	var sError string
	CSQL.QueryTb(C.SQL_TABLE.GameItem(), in, &listOut, &sError)

	if len(listOut) > 0 {

		C.InterfaceToStruct(listOut[0], &re)

	}

	return re
}

func GetUser(sSid string) (re C.UserData) {

	in := make(map[string]interface{})
	in["Sid"] = sSid

	listOut := []interface{}{}
	var sError string
	CSQL.QueryTb(C.SQL_TABLE.UserData(), in, &listOut, &sError)

	if len(listOut) > 0 {
		C.InterfaceToStruct(listOut[0], &re)

	}

	return re

}

func GetCustomer(sSid string) (re C.CustomerData) {
	in := make(map[string]interface{})
	in["Sid"] = sSid

	listOut := []interface{}{}
	var sError string
	CSQL.QueryTb(C.SQL_TABLE.CustomerData(), in, &listOut, &sError)

	if len(listOut) > 0 {
		C.InterfaceToStruct(listOut[0], &re)
	}

	return re

}
