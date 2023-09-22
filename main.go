package main

import (
	"fmt"
	"os"

	Adp "github.com/roger0816/adpGo/Adp"
	CSQL "github.com/roger0816/adpGo/CSql"
	C "github.com/roger0816/adpGo/Common"
	NETWORK "github.com/roger0816/adpGo/RpkNetwork"
)

func main() {

	args := os.Args[1:] // [1:] 可以跳過程序名稱
	fmt.Println("v1.1.0922")
	fmt.Printf("args : %v\n", args)
	listTen := "6005"
	dbIp := "172.104.112.34"

	if len(args) > 1 {
		dbIp = args[1]
	}

	if len(args) > 0 {
		listTen = args[0]
	}

	err := CSQL.OpenDb(dbIp, "3306", "adp", "roger", "Aa111111")

	if err != nil {
		fmt.Println("open db false")
		return
	}
	defer CSQL.CloseDb()

	runServer(listTen)

	select {}

}

func runServer(sPort string) {
	go NETWORK.StartTcpServer(sPort, Adp.AdpRecaller{})
	fmt.Println("server start " + sPort)
	iTmp := C.StringToInt64(sPort)

	sSubPort := C.Int64ToString(iTmp + 10)

	go NETWORK.StartTcpServer(sSubPort, Adp.AdpRecaller{})
	fmt.Println("server start " + sSubPort)

	sApiPort := C.Int64ToString(iTmp + 20)
	go NETWORK.StartApiServer(sApiPort)
	fmt.Println("api start " + sApiPort)

}

func test1() {
	var data NETWORK.CData

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
