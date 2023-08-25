package main

import (
	"fmt"

	CSQL "github.com/roger0816/adpGo/CSql"
	C "github.com/roger0816/adpGo/Common"
	NETWORK "github.com/roger0816/adpGo/RpkNetwork"
)

func main() {

	//test2()
	CSQL.OpenDb("172.104.112.34", "3306", "adp", "roger", "Aa111111")
	runServer()

	select {}
}

func runServer() {
	go NETWORK.StartTcpServer("6005")
	fmt.Println("server start 6005")
	go NETWORK.StartApiServer("6004")

	fmt.Println("server start 6004")
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
