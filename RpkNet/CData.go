package RpkNet

import (
	"encoding/json"
)

const END_DATA = "::ENDX::"

const (
	ACT_SEND   = 0
	ACT_RECALL = 1
)

type CData struct {
	Action    int           `json:"action"`
	User      string        `json:"user"`
	Msg       string        `json:"msg"`
	Ok        bool          `json:"ok"`
	State     int           `json:"status"`
	HeartBeat string        `json:"heartBeat"`
	Data      interface{}   `json:"data"`
	ListName  []interface{} `json:"listName"`
	ListData  []string      `json:"listData"`
	Trigger   string        `json:"trigger"`
	SendSync  interface{}   `json:"sendSync"`
	RecSync   interface{}   `json:"recSync"`
}

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
	return json.Marshal(d)
}

func (d *CData) DecodeJSON(jsonData []byte) error {
	return json.Unmarshal(jsonData, d)
}
