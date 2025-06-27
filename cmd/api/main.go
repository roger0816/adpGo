package main

import (
	"fmt"
	"os"

	C "adpGo/common"
	Adp "adpGo/internal/Adp"
	CSQL "adpGo/pkg/CSql"
	NETWORK "adpGo/pkg/RpkNetwork"
)

func main() {

	args := os.Args[1:] // [1:] 可以跳過程序名稱
	fmt.Println("v2.2.0124")
	fmt.Printf("args : %v\n", args)
	listTen := "6005"
	//test db
	dbIp := "172.104.112.34"
	//
	//dbIp := "172.104.117.7"

	dbName := "adp"

	Adp.OrderFuncIni()

	if len(args) > 1 {
		dbIp = args[1]
	}

	if len(args) > 2 {
		dbName = args[2]
	}

	if len(args) > 0 {
		listTen = args[0]
	}

	err := CSQL.OpenDb(dbIp, "3306", dbName, "roger", "Aa111111")

	if err != nil {
		fmt.Println("open db false")
		return
	}
	defer CSQL.CloseDb()

	runServer(listTen)

	select {}

}

func runServer(sPort string) {

	go Adp.RunApi(sPort)

}

func test1() {
	var data C.CData

	data.Action = 1000

	NETWORK.SendTcpData("127.0.0.1", "6005", data)
}

type testStruct1 struct {
	I    int    `json:"i" structs:"i"`
	Name string `json:"name" structs:"name"`
}

func test2() {
	dd1 := C.VariantMap{
		"i":    3,
		"name": "test",
		"id":   "a01",
		"tel":  "123456",
		"my":   "example",
	}

	//var myStruct testStruct1

	var myInterface interface{}

	//1 map to interface
	myInterface = dd1.ToInterface() // OR myInterface= dd1
	fmt.Println(myInterface)

	//2 map to struct

	// myStruct =dd1.(testStruct1)

	var dd2 map[string]interface{}
	C.InterFaceToMap(myInterface, &dd2)

	fmt.Println(dd2)

}

func test3() {

	var d C.OrderData

	d.Money = "AAAA;;BBBB"

	var list = d.GetList("Money")

	fmt.Printf("ddd0 %v \n", list)

	list = append(list, "CCCC")

	d.SetList("Money", list)

	fmt.Printf("ddd1 %v \n", d)

	d.AppendToList("Money", "DDDD")

}
