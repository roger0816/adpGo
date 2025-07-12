package RpkNetwork

import (
	C "adpGo/common"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

// Step 1: 定義一個 Recaller interface
type Recaller interface {
	ImplementRecall(data C.CData) C.CData
}

// 默認的實現
type DefaultRecaller struct{}

func (d DefaultRecaller) ImplementRecall(data C.CData) C.CData {
	// 你的默認操作
	return data // 返回修改後的data或原始data
}

func handleConnection(conn net.Conn, recaller Recaller) {
	defer conn.Close()

	var buffer bytes.Buffer
	tmp := make([]byte, 1024)

	for {

		for {

			// 设置读取超时
			conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

			n, err := conn.Read(tmp)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Client closed the connection")
				} else {
					fmt.Println("Error reading:", err)
				}

				return
			}

			buffer.Write(tmp[:n])

			//if bytes.Contains(buffer.Bytes(), []byte(footMaker)) {
			if IsPackageComplete(buffer.Bytes()) {
				break
			}
		}

		unBuff := UnPackage(buffer.Bytes())

		var bIsHeartBeat bool = strings.Contains(string(unBuff), "\"action\": 1,")
		var bMixQuery bool = strings.Contains(string(unBuff), "\"action\": 6031,")

		if !bIsHeartBeat {
			currentTime := time.Now()
			timeStr := currentTime.Format("15:04:05.999999")
			if !bMixQuery {
				fmt.Printf("[%s] server get data :\n %s \n\n", timeStr, string(unBuff))
			} else {
				fmt.Printf("[%s] server get data :mix request \n", timeStr)
			}
		}

		var decodedData C.CData
		err2 := decodedData.DecodeJSON(unBuff)
		if err2 != nil {
			fmt.Println("Error decoding JSON:", err2)
			return
		}

		//re := ImplementRecall(decodedData)

		re := recaller.ImplementRecall(decodedData)
		reEncoded, err3 := re.EncodeJSON()

		if err3 != nil {
			fmt.Println("Error encoding JSON:", err3)
			return
		}

		recall := Package(reEncoded)

		_, err4 := conn.Write(recall)

		if err4 != nil {
			fmt.Println("Error writing:", err4)
			return
		}

		if !bIsHeartBeat {

			currentTime2 := time.Now()
			timeStr := currentTime2.Format("15:04:05.999999")
			if !bMixQuery {
				fmt.Printf("[%s]Received: %s\n", timeStr, reEncoded)
			} else {
				fmt.Printf("[%s]Received: mixdata\n", timeStr)
			}
		}

		buffer.Reset() // 清空 buffer，等待下一包
	}

}

/*
func handleConnection(conn net.Conn, recaller Recaller) {
	defer conn.Close()

	var buffer bytes.Buffer
	tmp := make([]byte, 1024)

	for {

		// 设置读取超时
		conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

		n, err := conn.Read(tmp)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed the connection")
			} else {
				fmt.Println("Error reading:", err)
			}

			return
		}

		buffer.Write(tmp[:n])

		//if bytes.Contains(buffer.Bytes(), []byte(footMaker)) {
		if IsPackageComplete(buffer.Bytes()) {
			break
		}
	}

	unBuff := UnPackage(buffer.Bytes())

	var bIsHeartBeat bool = strings.Contains(string(unBuff), "\"action\": 1,")
	var bMixQuery bool = strings.Contains(string(unBuff), "\"action\": 6031,")

	if !bIsHeartBeat {
		currentTime := time.Now()
		timeStr := currentTime.Format("15:04:05.999999")
		if !bMixQuery {
			fmt.Printf("[%s] server get data :\n %s \n\n", timeStr, string(unBuff))
		} else {
			fmt.Printf("[%s] server get data :mix request \n", timeStr)
		}
	}

	var decodedData C.CData
	err2 := decodedData.DecodeJSON(unBuff)
	if err2 != nil {
		fmt.Println("Error decoding JSON:", err2)
		return
	}

	//re := ImplementRecall(decodedData)

	re := recaller.ImplementRecall(decodedData)
	reEncoded, err3 := re.EncodeJSON()

	if err3 != nil {
		fmt.Println("Error encoding JSON:", err3)
		return
	}

	recall := Package(reEncoded)

	_, err4 := conn.Write(recall)

	if err4 != nil {
		fmt.Println("Error writing:", err4)
		return
	}

	if !bIsHeartBeat {

		currentTime2 := time.Now()
		timeStr := currentTime2.Format("15:04:05.999999")
		if !bMixQuery {
			fmt.Printf("[%s]Received: %s\n", timeStr, reEncoded)
		} else {
			fmt.Printf("[%s]Received: mixdata\n", timeStr)
		}
	}

}
*/

// 可以提供一個默認的實現或者允許外部提供
func StartTcpServer(port string, recaller Recaller) {

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		if recaller == nil {
			recaller = DefaultRecaller{}
		}
		go handleConnection(conn, recaller)
	}
}

func SendTcp(ip, port, message string) (string, error) {
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

func SendTcpData(ip, port string, data C.CData) (C.CData, error) {

	fmt.Println("send tcp data")

	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return C.CData{}, fmt.Errorf("Error connecting: %v", err)
	}
	defer conn.Close()

	encodedData, err := data.EncodeJSON()
	if err != nil {
		return C.CData{}, fmt.Errorf("Error encoding data: %v", err)
	}
	fmt.Printf("send tcp : %x", encodedData)
	_, err = conn.Write(encodedData)
	if err != nil {
		return C.CData{}, fmt.Errorf("Error sending data: %v", err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return C.CData{}, fmt.Errorf("Error receiving data: %v", err)
	}

	receivedData := buffer[:n]
	var response C.CData
	err = response.DecodeJSON(receivedData)
	if err != nil {
		return C.CData{}, fmt.Errorf("Error decoding response: %v", err)
	}

	return response, nil
}
