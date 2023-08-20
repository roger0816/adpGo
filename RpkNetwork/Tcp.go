package RpkNetwork

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"github.com/roger0816/adpGo/CSql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	headerMaker = "@AKA404@"
	footMaker   = "@404AKA@"
	endMarker   = "::ENDX::"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var buffer bytes.Buffer
	tmp := make([]byte, 1024)

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		buffer.Write(tmp[:n])

		if bytes.Contains(buffer.Bytes(), []byte(footMaker)) {
			break
		}
	}

	quotedData := strconv.QuoteToASCII(string(buffer.Bytes()))
	quotedData = strings.Replace(quotedData, headerMaker, "", -1)
	quotedData = strings.Replace(quotedData, footMaker, "", -1)
	quotedData = strings.Replace(quotedData, endMarker, "", -1)
	quotedData = strings.Replace(quotedData, "\\x00", "", -1)
	quotedData = strings.Replace(quotedData, "\\x01_", "", -1)
	quotedData = strings.Replace(quotedData, " ", "", -1)
	quotedData = strings.Replace(quotedData, "\\n", "", -1)


	fmt.Printf("server get data :\n %s \n\n",quotedData)

	if len(quotedData) > 2 {

		if quotedData[0] == '"' {
			quotedData = quotedData[1:]

		}

		if quotedData[len(quotedData)-1] == '"' {
			quotedData = quotedData[:len(quotedData)-1]
		}

	}

	quotedData = strings.Replace(quotedData, "\\", "", -1)

	tmpByteArray := []byte(quotedData)

	var decodedData CData
	err2 := decodedData.DecodeJSON(tmpByteArray)
	if err2 != nil {
		fmt.Println("Error decoding JSON:", err2)
		return
	}

	re := Query(decodedData)
	reEncoded, err3 := re.EncodeJSON()

	if err3 != nil {
		fmt.Println("Error encoding JSON:", err3)
		return
	}

	_, err4 := conn.Write(reEncoded)

	if err4 != nil {
		fmt.Println("Error writing:", err4)
		return
	}

	fmt.Printf("Decoded Data: %+v\n", decodedData)

	fmt.Printf("Received: %s\n", quotedData)
}


func TcpListen(port string) {

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Listening on port", port)


	 CSql.OpenDb("172.104.112.34", "3306", "adp", "roger", "Aa111111")




	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}


func SendTcp(ip ,port ,message string)(string ,error ){
	conn, err := net.Dial("tcp", ip+":"+port)

	if err != nil {
		return "", fmt.Errorf("Error connecting: %v", err)
	}
	defer conn.Close()

	// 发送数据
	_, err = conn.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("Error sending data: %v", err)
	}

	// 接收数据
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("Error receiving data: %v", err)
	}

	response := string(buffer[:n])
	return response, nil

}


func SendTcpData(ip,port string, data CData) (CData, error) {

	fmt.Println("send tcp data")

	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return CData{}, fmt.Errorf("Error connecting: %v", err)
	}
	defer conn.Close()

	encodedData, err := data.EncodeJSON()
	if err != nil {
		return CData{}, fmt.Errorf("Error encoding data: %v", err)
	}
	fmt.Printf("send tcp : %x",encodedData)
	_, err = conn.Write(encodedData)
	if err != nil {
		return CData{}, fmt.Errorf("Error sending data: %v", err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return CData{}, fmt.Errorf("Error receiving data: %v", err)
	}

	receivedData := buffer[:n]
	var response CData
	err = response.DecodeJSON(receivedData)
	if err != nil {
		return CData{}, fmt.Errorf("Error decoding response: %v", err)
	}

	return response, nil
}