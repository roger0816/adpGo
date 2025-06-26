package RpkNetwork

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func headerByte() []byte {
	return []byte{'@', 'A', 'K', 'A', '4', '0', '4', '@'}
}

func footByte() []byte {
	return []byte{'@', '4', '0', '4', 'A', 'K', 'A', '@'}
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

func IsPackageComplete(oriData []byte) bool {

	if bytes.Contains(oriData, []byte(headerByte())) &&
		bytes.Contains(oriData, []byte(footByte())) {
		return true
	}

	return false

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
