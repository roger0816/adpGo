package RpkNetwork

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

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

	unBuff := UnPackage(buffer.Bytes())

	var  bIsHeartBeat bool =strings.Contains(string(unBuff),"\"action\": 1,")
	var  bMixQuery bool =strings.Contains(string(unBuff),"\"action\": 6031,")



	if !bIsHeartBeat && !bMixQuery{
		currentTime := time.Now()
		timeStr :=currentTime.Format("15:04:05.999999")
	
	fmt.Printf("[%s] server get data :\n %s \n\n",timeStr, string(unBuff))
	}

	var decodedData CData
	err2 := decodedData.DecodeJSON(unBuff)
	if err2 != nil {
		fmt.Println("Error decoding JSON:", err2)
		return
	}

	re := ImplementRecall(decodedData)
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

	if !bIsHeartBeat && !bMixQuery{

		currentTime2 := time.Now()
		timeStr :=currentTime2.Format("15:04:05.999999")

	fmt.Printf("[%s]Received: %s\n",timeStr, reEncoded)
	}


}

func StartTcpServer(port string) {

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
		go handleConnection(conn)
	}
}



//
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

func SendTcpData(ip, port string, data CData) (CData, error) {

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
	fmt.Printf("send tcp : %x", encodedData)
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


func Uint32ToByteArray(somevalue uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, somevalue)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return buf.Bytes()
}

func ByteArrayToUint32(data []byte) uint32 {
	//return binary.LittleEndian.Uint32(data)
	return binary.BigEndian.Uint32(data)
}

func headerByte() []byte {
	return []byte{'@', 'A', 'K', 'A', '4', '0', '4', '@'}
}

func footByte() []byte {
	return []byte{'@', '4', '0', '4', 'A', 'K', 'A', '@'}
}

func Package(data []byte) []byte {
	var pdata []byte
	//data = append(data, []byte(endMarker)...)
	pdata = append(pdata, headerByte()...)
	dataLength := uint32(len(data) + 20)
	dataLengthBytes := Uint32ToByteArray(dataLength)

	pdata = append(pdata, dataLengthBytes...)
	pdata = append(pdata, data...)
	pdata = append(pdata, footByte()...)
	return pdata
}

func UnPackage(data []byte) []byte {
	hb := headerByte()
	ft := footByte()

	sp := bytes.Index(data, hb)

	if sp < 0 {
		// fmt.Println("unPackage:No header", string(p.m_data))
		return nil
	}

	ep := bytes.Index(data, ft)

	if ep < 0 {
		// fmt.Println("unPackage:No footer", string(p.m_data))
		return nil
	}

	datavlen := ByteArrayToUint32(data[sp+8 : sp+12])

	rdata := data[sp+12 : sp+12+int(datavlen)-20]

	// fmt.Println("unPackage", string(rdata))

	return rdata
}
