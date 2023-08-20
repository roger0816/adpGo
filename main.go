package main

import (
	"fmt"

	"github.com/roger0816/adpGo/RpkNetwork"
)

func main() {

	/*

		data := CNetwork.CData{
			Action:    CNetwork.ACT_SEND,
			User:      "user123",
			Msg:       "Hello, World!",
			Ok:        true,
			State:     0,
			HeartBeat: "heartbeat",
			Data:      nil,
			ListName:  nil,
			ListData:  nil,
			Trigger:   "trigger123",
			SendSync:  nil,
			RecSync:   nil,
		}



		encodedData, err := data.EncodeJSON()
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		fmt.Println("Encoded JSON:", string(encodedData))

		var decodedData CNetwork.CData
		err = decodedData.DecodeJSON(encodedData)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		fmt.Printf("Decoded Data: %+v\n", decodedData)
	*/
	RpkNetwork.TcpListen("6005")

	fmt.Print("server start")

	/*
		var input string
		fmt.Print("Enter something: ")
		fmt.Scanln(&input)
		fmt.Println("You entered:", input)
	*/
	//test1()

}

func test1() {
	var data CNetwork.CData

	data.Action = 1000

	CNetwork.SendTcpData("127.0.0.1", "6005", data)
}
