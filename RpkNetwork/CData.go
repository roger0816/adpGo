package RpkNetwork

import (
	"encoding/json"
	"strings"
	//"fmt"
)

const END_DATA = "::ENDX::"

const (
	ACT_SEND   = 0
	ACT_RECALL = 1
	ACT_LOCAL =3
)

type CData struct {
	Action    int           `json:"action"`
	User      string        `json:"user"`
	Msg       string        `json:"msg"`
	Ok        bool          `json:"ok"`
	State     int           `json:"status"`
	HeartBeat string        `json:"heartBeat"`
	Data      map[string]interface{}   `json:"data"`
	ListName  []interface{} `json:"listName"`
	ListData  []interface{} `json:"listData"`
	Trigger   string        `json:"trigger"`
	SendSync  interface{}   `json:"sendSync"`
	RecSync   interface{}   `json:"recSync"`
}

// func (d *CData) DataToMap() map[string]interface{} {
// 	var t map[string]interface{}
// 	if dataMap, ok := d.Data.(map[string]interface{}); ok {
// 		t = dataMap

// 	} else {
// 		// 如果无法转换为map[string]interface{}，进行相应的错误处理
// 		// fmt.Println("Failed to convert re.Data to map[string]interface{}")
// 	}

// 	return t
// }

func NewData() CData {
	return CData{
		Action:    ACT_SEND,
		User:      "",
		Msg:       "",
		Ok:        false,
		State:     0,
		HeartBeat: "",
		Data:      nil,
		ListName:  nil,
		ListData:  nil,
		Trigger:   "",
		SendSync:  nil,
		RecSync:   nil,
	}
}

func (d CData) EncodeJSON() ([]byte, error) {
	data , err :=json.Marshal(d)
	if err != nil {
		return nil ,err
	}
	data = append(data, []byte(endMarker)...)
	return data,nil
}

func (d *CData) DecodeJSON(jsonData []byte) error {
	quotedData := string(jsonData)
    quotedData = strings.Replace(quotedData, endMarker, "", -1)
	buffByte:=[]byte(quotedData)
	return json.Unmarshal(buffByte, d)
}
