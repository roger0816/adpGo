package RpkNetwork

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	CSQL "github.com/roger0816/adpGo/CSql"
	C "github.com/roger0816/adpGo/Common"
)

var LoginUser = make(map[string]interface{})

func ImplementRecall(data CData) CData {

	var reData C.VariantMap = make(map[string]interface{})
	reList := []interface{}{}
	iAction := data.Action
	var error, sOkMsg string
	var bOk bool = false

	var tmpIn C.VariantMap = make(map[string]interface{})

	var tmpMap C.VariantMap = make(map[string]interface{})
	//var tmpList = []interface{}{}
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
		fmt.Println("login")
		//conditions:=[string]interface
		//CSql.QueryTb(tableName, conditions, listOut, sError)
		bOk = false
		error = "連線失敗"

		if len(data.ListData) >= 2 {

			in := map[string]interface{}{
				"Id":       data.ListData[0],
				"Password": data.ListData[1],
			}

			tmpB := CSQL.QueryTb(C.SQL_TABLE.UserData(), in, &reList, &error)

			if tmpB {
				error = "帳密錯誤"
			}

			if len(reList) > 0 {

				strBytes := []byte(C.TimeUtc8Str())
				hash := md5.Sum(strBytes)
				sSession := hex.EncodeToString(hash[:])
				reList[0].(map[string]interface{})["Session"] = sSession
				fmt.Println(sSession)
				C.InterFaceToMap(reList[0], reData.Origin())

				var UserSid string
				if sid, ok := reList[0].(map[string]interface{})["Sid"].(string); ok {
					UserSid = sid
				} else {
					// 处理无法转换为字符串的情况
					fmt.Println("Failed to convert Sid to string")
				}

				LoginUser[UserSid] = sSession
				bOk = true
				sOkMsg = "登入成功"

				fmt.Println("loging msg : " + sOkMsg)
			}

		}

		//
	case iAction == C.SET_VALUE:

		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.Settings(), Data, &error, true)

	//QUERY============

	case iAction == C.GET_VALUE:

		bOk = CSQL.QueryTb(C.SQL_TABLE.Settings(), Data, &reList, &error)
		if bOk && len(reList) > 0 {
			C.InterFaceToMap(reList[0], reData.Origin())
		}

	case iAction == C.QUERY_USER:

		bOk = CSQL.QueryTb(C.SQL_TABLE.UserData(), Data, &reList, &error)

	case iAction == C.QUERY_CUSTOMER:
		bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerData(), Data, &reList, &error)

	case iAction == C.QUERY_CUSTOMER_GAME_INFO:
		bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerGameInfo(), Data, &reList, &error)

	case iAction == C.QUERY_CUSTOMER_COST:
		bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerCost(), Data, &reList, &error)

	case iAction == C.QUERY_GAME_LIST:
		bOk = CSQL.QueryTb(C.SQL_TABLE.GameList(), Data, &reList, &error)

	case iAction == C.QUERY_GAME_RATE:

		bOk = CSQL.QueryTb(C.SQL_TABLE.GameRate(), Data, &reList, &error)

		if bOk && len(reList) < 1 {
			tmpOut := make([]interface{}, 0)
			tmp := map[string]interface{}{
				"sid": Data["GameSid"],
			}

			bOk = CSQL.QueryTb(C.SQL_TABLE.GameList(), tmp, &tmpOut, &error)

			if bOk && len(tmpOut) > 0 {
				//to do
			}

		}

	case iAction == C.QUERY_GAME_ITEM:
		bOk = CSQL.QueryTb(C.SQL_TABLE.GameItem(), Data, &reList, &error)

	case iAction == C.QUERY_BULLETIN:
		bOk = CSQL.QueryTb(C.SQL_TABLE.Bulletin(), Data, &reList, &error)

	case iAction == C.QUERY_CUSTOM_CLASS:
		in := make(map[string]interface{})
		in["Type"] = "group"
		bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerClass(), in, &reList, &error)

	case iAction == C.QUERY_CUSTOM_DEBIT:
		in := make(map[string]interface{})
		in["Type"] = "debit"
		bOk = CSQL.QueryTb(C.SQL_TABLE.CustomerClass(), in, &reList, &error)

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

		bOk = CSQL.QueryTb(C.SQL_TABLE.FactoryClass(), d, &reList, &error)

		fac := C.DataFactory{
			Sid:  999,
			Id:   "未",
			Name: "未分配",
		}
		reList = append(reList, fac)

	case iAction == C.QUERY_GROUP:
		bOk = CSQL.QueryTb(C.SQL_TABLE.GroupData(), Data, &reList, &error)

	case iAction == C.QUERY_PAY_TYPE:
		bOk = CSQL.QueryTb(C.SQL_TABLE.PayType(), Data, &reList, &error)

	case iAction == C.QUERY_ORDER:
		bOk = CSQL.QueryTb(C.SQL_TABLE.OrderData(), Data, &reList, &error)

	case iAction == C.QUERY_BOUNS:
		bOk = CSQL.QueryTb(C.SQL_TABLE.UserBonus(), Data, &reList, &error)

	case iAction == C.QUERY_SCHEDULE:
		bOk = CSQL.QueryTb(C.SQL_TABLE.Schedule(), Data, &reList, &error)

	case iAction == C.QUERY_EXCHANGE:
		bOk = CSQL.QueryTb(C.SQL_TABLE.ExchangeRate(), Data, &reList, &error)

	case iAction == C.QUERY_PRIMERATE:
		bOk = CSQL.QueryTb(C.SQL_TABLE.PrimeCostRate(), Data, &reList, &error)

	case iAction == C.QUERY_PIC:
		bOk = CSQL.QueryTb(C.SQL_TABLE.PicData(), Data, &reList, &error)

	case iAction == C.QUERY_ITEM_COUNT:
		bOk = CSQL.QueryTb(C.SQL_TABLE.GameItemCount(), Data, &reList, &error)

	case iAction == C.QUERY_DEBIT_CLASS:
		tmpIn := make(map[string]interface{})
		tmpIn["ASC"] = "Sort"

		bOk = CSQL.QueryTb(C.SQL_TABLE.DebitClass(), tmpIn, &reList, &error)

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

			var tmp CData

			if err := tmp.DecodeJSON(buff); err != nil {
				continue
			}

			tmp.State = ACT_LOCAL
			tmpRe := ImplementRecall(tmp)

			if tmp2, err := tmpRe.EncodeJSON(); err == nil {
				//fmt.Printf("%d encode ok ", api)
				reData[key] = tmp2
			}

		}

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

		bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.UserData(), C.StructToMap(user), &error, false)
		sOkMsg = "新增成功"

	case iAction == C.EDIT_USER:
		in := make(map[string]interface{})
		in["Sid"] = Data["Sid"]
		bOk = CSQL.UpdateTb(C.SQL_TABLE.UserData(), in, Data, &error)
		sOkMsg = "修改成功"

	case iAction == C.DEL_USER:
		bOk = CSQL.DelFromTb(C.SQL_TABLE.UserData(), Data, &error)
		sOkMsg = "刪除成功"

	case iAction == C.ADD_CUSTOMER:

		bOk, _, tmpMap = CSQL.InsertTb(C.SQL_TABLE.CustomerData(), Data, &error, false)

		sOkMsg = "客戶名稱:" + tmpMap.String("Name") + "\n新增成功"

	case iAction == C.EDIT_CUSTOMER:

		tmpIn["Sid"] = Data["Sid"]
		bOk = CSQL.UpdateTb(C.SQL_TABLE.CustomerData(), tmpIn, Data, &error)
		sOkMsg = "修改成功"

	case iAction == C.DEL_CUSTOMER:
		in := make(map[string]interface{})
		in["Sid"] = Data["Sid"]
		bOk = CSQL.DelFromTb(C.SQL_TABLE.CustomerData(), in, &error)
		if bOk {
			CSQL.DelFromTb(C.SQL_TABLE.CustomerMoney(), in, &error)

			in = make(map[string]interface{})
			in["CustomerSid"] = Data["Sid"]
			CSQL.DelFromTb(C.SQL_TABLE.CustomerGameInfo(), in, &error)

			CSQL.DelFromTb(C.SQL_TABLE.CustomerCost(), in, &error)
		}

		sOkMsg = "刪除成功"

	case iAction == C.REPLACE_GAME_INFO:

		for i := 0; i < len(data.ListData); i++ {
			bTmp, _,_ := CSQL.InsertTb(C.SQL_TABLE.CustomerGameInfo(), data.ListData[i].(map[string]interface{}),&error,true)
		
			if  !bTmp {
				bOk = false
			}
		}

		sOkMsg = "客戶遊戲資料修改完成"
		
	case iAction == C.DEL_GAME_INFO:
		for i := 0; i < len(data.ListData); i++ {
		   C.InterFaceToMap(data.ListData[i],tmpMap.Origin()) 
		   	tmpIn["Sid"] = tmpMap.String("Sid")
			bTmp:= CSQL.DelFromTb(C.SQL_TABLE.CustomerGameInfo(), tmpIn,&error)
		
			if  !bTmp {
				bOk = false
			}
		}

		sOkMsg = "客戶遊戲資料修改完成"


	case iAction == C.EDIT_GAME_LIST:

		in := map[string]interface{}{
			"Sid": Data["Sid"],
		}

		var game C.DataGameList

		C.SetData(Data, &game)

		//tmpRate := strconv.FormatFloat(game.GameRate, 'f', 2, 64) // 'f' 表示格式，2 表示小數點後的位數，64 表示它是 float64

		gameRate := map[string]interface{}{
			"GameSid":  game.Id,
			"GameName": game.Name,
			"UserSid":  game.UserSid,
			"Rate":     game.GameRate,
		}

		bOk = CSQL.UpdateTb(C.SQL_TABLE.GameList(), in, C.GetData(&game), &error)

		if bOk {
			bOk, _, _ = CSQL.InsertTb(C.SQL_TABLE.GameRate(), gameRate, &error, true)
		}
		sOkMsg = "修改完成"

	case iAction == C.EDIT_GAME_ITEM:
		fmt.Println("Edit game item:")

		bHasList := len(data.ListData) > 0

		bOk = true

		if !bHasList {
			d := make(map[string]interface{})
			d["Sid"] = Data["Sid"]

			bOk = CSQL.UpdateTb(C.SQL_TABLE.GameItem(), d, Data, &error)

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
				bOk = CSQL.BatchUpdateTb(C.SQL_TABLE.GameItem(), listConditions, listMap, &error)

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

					// 假設 CSQL.UpdateTb 是一個用來更新數據表的Go函數
					b := CSQL.UpdateTb("GameItem", tmp, tmpMap, &error)
					if !b {
						bOk = false
					}
				}

			}

		}

		sOkMsg = "修改完成"
		fmt.Println(sOkMsg)

	default:
		// 未知操作，可以进行相应的处理
		fmt.Printf("unknown action : %d \n", iAction)

	}

	var re CData
	re.Action = data.Action
	re.State = ACT_RECALL

	re.Ok = bOk

	if re.Ok {
		re.Msg = sOkMsg
	} else {
		re.Msg = error
	}

	re.Data = reData

	re.ListData = reList
	return re
}
