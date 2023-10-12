package Common

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type UserData struct {
	Sid        int    `json:"Sid" structs:"Sid"`
	Id         string `json:"Id" structs:"Id"`
	Name       string `json:"Name" structs:"Name"`
	UpdateTime string `json:"UpdateTime" structs:"UpdateTime"`
	Password   string `json:"Password" structs:"Password"`
	Cid        string `json:"Cid" structs:"Cid"`
	Lv         int    `json:"Lv" structs:"Lv"`
	ParentId   string `json:"ParentId" structs:"ParentId"`
	StartDay   string `json:"StartDay" structs:"StartDay"`
	BirthDay   string `json:"BirthDay" structs:"BirthDay"`
	Tel        string `json:"Tel" structs:"Tel"`
	Email      string `json:"Email" structs:"Email"`
	Note1      string `json:"Note1" structs:"Note1"`
	Note2      string `json:"Note2" structs:"Note2"`
	Note3      string `json:"Note3" structs:"Note3"`
	CreateTime string `json:"CreateTime" structs:"CreateTime"`
}

type DataUserBonus struct {
	Sid        int    `json:"Sid" structs:"Sid"`
	UserSid    string `json:"UserSid" structs:"UserSid"`
	OrderSid   string `json:"OrderSid" structs:"OrderSid"`
	Bonus      string `json:"Bonus" structs:"Bonus"`
	Pay        string `json:"Pay" structs:"Pay"`
	AddBonus   string `json:"AddBonus" structs:"AddBonus"`
	AddPay     string `json:"AddPay" structs:"AddPay"`
	UpdateTime string `json:"UpdateTime" structs:"UpdateTime"`
}

type CustomerData struct {
	Sid        int    `json:"Sid" structs:"Sid"`
	Id         string `json:"Id" structs:"Id"`
	Name       string `json:"Name" structs:"Name"`
	UpdateTime string `json:"UpdateTime" structs:"UpdateTime"`
	Vip        string `json:"Vip" structs:"Vip"`
	Class      string `json:"Class" structs:"Class"`
	Money      string `json:"Money" structs:"Money"`
	Line       string `json:"Line" structs:"Line"`
	Currency   string `json:"Currency" structs:"Currency"`
	PayType    string `json:"PayType" structs:"PayType"`
	PayInfo    string `json:"PayInfo" structs:"PayInfo"`
	UserSid    string `json:"UserSid" structs:"UserSid"`
	Num5       string `json:"Num5" structs:"Num5"`
	Note1      string `json:"Note1" structs:"Note1"`
	Note2      string `json:"Note2" structs:"Note2"`
}

type CustomerMoney struct {
	Sid        int    `json:"Sid" structs:"Sid"`
	Id         string `json:"Id" structs:"Id"`
	Name       string `json:"Name" structs:"Name"`
	UpdateTime string `json:"UpdateTime" structs:"UpdateTime"`
	Money      string `json:"Money" structs:"Money"`
	Currency   string `json:"Currency" structs:"Currency"`
}

type CustomerCost struct {
	Sid            string `json:"Sid" structs:"Sid"`
	CustomerSid    string `json:"CustomerSid" structs:"CustomerSid"`
	OrderId        string `json:"OrderId" structs:"OrderId"`
	Currency       string `json:"Currency" structs:"Currency"`
	ChangeValue    string `json:"ChangeValue" structs:"ChangeValue"`
	OriginCurrency string `json:"OriginCurrency" structs:"OriginCurrency"`
	OriginValue    string `json:"OriginValue" structs:"OriginValue"`
	DebitSid       string `json:"DebitSid" structs:"DebitSid"`
	DebitNote      string `json:"DebitNote" structs:"DebitNote"`
	Rate           string `json:"Rate" structs:"Rate"`
	AddRate        string `json:"AddRate" structs:"AddRate"`
	IsAddCost      bool   `json:"IsAddCost" structs:"IsAddCost"`
	Total          string `json:"Total" structs:"Total"`
	UserSid        string `json:"UserSid" structs:"UserSid"`
	OrderTime      string `json:"OrderTime" structs:"OrderTime"`
	UpdateTime     string `json:"UpdateTime" structs:"UpdateTime"`
	Pic0           string `json:"Pic0" structs:"Pic0"`
	Pic1           string `json:"Pic1" structs:"Pic1"`
	Note0          string `json:"Note0" structs:"Note0"`
	Note1          string `json:"Note1" structs:"Note1"`
}

type DataCustomerClass struct {
	Sid        int    `json:"Sid" structs:"Sid"`
	Id         string `json:"Id" structs:"Id"`
	Name       string `json:"Name" structs:"Name"`
	UpdateTime string `json:"UpdateTime" structs:"UpdateTime"`
	Sort       int    `json:"Sort" structs:"Sort"`
	Type       string `json:"Type" structs:"Type"`
	Note1      string `json:"Note1" structs:"Note1"`
	Note2      string `json:"Note2" structs:"Note"`
}
type DebitClass struct {
	Sid        int    `json:"Sid" structs:"Sid"`
	Id         string `json:"Id" structs:"Id"`
	Name       string `json:"Name" structs:"Name"`
	UpdateTime string `json:"UpdateTime" structs:"UpdateTime"`
	Sort       int    `json:"Sort" structs:"Sort"`
	Currency   string `json:"Currency" structs:"Currency"`
	Note1      string `json:"Note1" structs:"Note1"`
	Note2      string `json:"Note2" structs:"Note2"`
}

type GroupData struct {
	Sid        int    `json:"Sid" structs:"Sid"`
	Id         string `json:"Id" structs:"Id"`
	Name       string `json:"Name" structs:"Name"`
	UpdateTime string `json:"UpdateTime" structs:"UpdateTime"`
	Type       string `json:"Type" structs:"Type"`
	Value      string `json:"Value" structs:"Value"`
	Blob       []byte `json:"Blob" structs:"Blob"`
	Note1      string `json:"Note1" structs:"Note1"`
	Note2      string `json:"Note2" structs:"Note2"`
	NoteBlob   []byte `json:"NoteBlob" structs:"NoteBlob"`
}

type CustomerGameInfo struct {
	Sid           int    `json:"Sid" structs:"Sid"`
	Id            string `json:"Id" structs:"Id"`
	Name          string `json:"Name" structs:"Name"`
	UpdateTime    string `json:"UpdateTime" structs:"UpdateTime"`
	CustomerSid   string `json:"CustomerSid" structs:"CustomerSid"`
	CustomerId    string `json:"CustomerId" structs:"CustomerId"`
	GameSid       string `json:"GameSid" structs:"GameSid"`
	LoginType     string `json:"LoginType" structs:"LoginType"`
	LoginAccount  string `json:"LoginAccount" structs:"LoginAccount"`
	LoginPassword string `json:"LoginPassword" structs:"LoginPassword"`
	ServerName    string `json:"ServerName" structs:"ServerName"`
	Characters    string `json:"Characters" structs:"Characters"`
	LastTime      string `json:"LastTime" structs:"LastTime"`
	Note1         string `json:"Note1" structs:"Note1"`
}

type OrderData struct {
	Sid              int    `json:"Sid" structs:"Sid"`
	Id               string `json:"Id" structs:"Id"`
	Name             string `json:"Name" structs:"Name"`
	UpdateTime       string `json:"UpdateTime" structs:"UpdateTime"`
	CustomerSid      string `json:"CustomerSid" structs:"CustomerSid"`
	UiRecord         string `json:"UiRecord" structs:"UiRecord"`
	Step             string `json:"Step" structs:"Step"`
	StepTime         string `json:"StepTime" structs:"StepTime"`
	User             string `json:"User" structs:"User"`
	Owner            string `json:"Owner" structs:"Owner"`
	PaddingUser      string `json:"PaddingUser" structs:"PaddingUser"`
	GameSid          string `json:"GameSid" structs:"GameSid"`
	Item             string `json:"Item" structs:"Item"`
	Cost             string `json:"Cost" structs:"Cost"`
	Bouns            string `json:"Bouns" structs:"Bouns"`
	PayType          string `json:"PayType" structs:"PayType"`
	CanSelectPayType string `json:"CanSelectPayType" structs:"CanSelectPayType"`
	GameRate         string `json:"GameRate" structs:"GameRate"`
	ExRateSid        string `json:"ExRateSid" structs:"ExRateSid"`
	PrimeRateSid     string `json:"PrimeRateSid" structs:"PrimeRateSid"`
	Money            string `json:"Money" structs:"Money"`
	Note0            string `json:"Note0" structs:"Note0"`
	ListCost         string `json:"ListCost" structs:"ListCost"`
	ListBouns        string `json:"ListBouns" structs:"ListBouns"`
	ItemInfo         string `json:"ItemInfo" structs:"ItemInfo"`
	Note1            string `json:"Note1" structs:"Note1"`
	Note2            string `json:"Note2" structs:"Note2"`
	Note3            string `json:"Note3" structs:"Note3"`
	Note4            string `json:"Note4" structs:"Note4"`
	Note5            string `json:"Note5" structs:"Note5"`
	Pic0             string `json:"Pic0" structs:"Pic0"`
	Pic1             string `json:"Pic1" structs:"Pic1"`
	OrderDate        string `json:"OrderDate" structs:"OrderDate"`
	OrderTime        string `json:"OrderTime" structs:"OrderTime"`
	CustomerName     string `json:"CustomerName" structs:"CustomerName"`
	Currency         string `json:"Currency" structs:"Currency"`
}

func (d *OrderData) SetList(sKey string, listData []string) {
	//因為這些用;; 間隔

	for len(listData) < 6 {
		listData = append(listData, "")
	}

	var str string
	if sKey == "StepTime" || sKey == "User" {
		str = strings.Join(listData, ",")
	} else {
		str = strings.Join(listData, ";;")
	}
	// // 雖然可以用反射，但不是每個都可以切的，想明碼
	// rv := reflect.ValueOf(d).Elem()
	// fieldValue := rv.FieldByName(sKey)
	// if fieldValue.IsValid() && fieldValue.CanSet() {
	// 	fieldValue.SetString(str)
	// }

	switch sKey {
	case "StepTime":
		d.StepTime = str
	case "User":
		d.User = str
	case "Item":
		d.Item = str
	case "Money":
		d.Money = str
	case "Note0":
		d.Note0 = str
	case "ListCost":
		d.ListCost = str
	case "ListBouns":
		d.ListBouns = str
	case "ItemInfo":
		d.ItemInfo = str
	case "CanSelectPayType":
		d.CanSelectPayType = str

	default:

	}

}

func (d *OrderData) GetList(sKey string) []string {
	//因為這些用;; 間隔
	var reList []string
	switch sKey {
	case "StepTime":
		reList = strings.Split(d.StepTime, ",")
	case "User":
		reList = strings.Split(d.User, ",")

	case "Item":
		reList = strings.Split(d.Item, ";;")
	case "Money":
		reList = strings.Split(d.Money, ";;")
	case "Note0":
		reList = strings.Split(d.Note0, ";;")
	case "ListCost":
		reList = strings.Split(d.ListCost, ";;")
	case "ListBouns":
		reList = strings.Split(d.ListBouns, ";;")
	case "ItemInfo":
		reList = strings.Split(d.ItemInfo, ";;")
	case "CanSelectPayType":
		reList = strings.Split(d.CanSelectPayType, ";;")

	}

	return reList
}

func (d *OrderData) AppendToList(sKey string, value string) {
	list := d.GetList(sKey)
	list = append(list, value)
	d.SetList(sKey, list)
}

type DataFactory struct {
	Sid        int      `json:"Sid" structs:"Sid"`
	Id         string   `json:"Id" structs:"Id"`
	Name       string   `json:"Name" structs:"Name"`
	UpdateTime string   `json:"UpdateTime" structs:"UpdateTime"`
	Currency   string   `json:"Currency" structs:"Currency"`
	PayTypdSid []string `json:"PayTypdSid" structs:"PayTypdSid"`
}

type DataGameList struct {
	Sid        int    `json:"Sid" structs:"Sid"`
	Id         string `json:"Id" structs:"Id"`
	Name       string `json:"Name" structs:"Name"`
	UpdateTime string `json:"UpdateTime" structs:"UpdateTime"`
	Enable     bool   `json:"Enable" structs:"Enable"`
	GameRate   string `json:"GameRate" structs:"GameRate"`
	UserSid    string `json:"UserSid" structs:"UserSid"`
	SellNote   string `json:"SellNote" structs:"SellNote"`
}

type DataGameRate struct {
	Sid        int    `json:"Sid" structs:"Sid"`
	GameSid    string `json:"GameSid" structs:"GameSid"`
	GameName   string `json:"GameName" structs:"GameName"`
	Rate       string `json:"Rate" structs:"Rate"`
	UserSid    string `json:"UserSid" structs:"UserSid"`
	UpdateTime string `json:"UpdateTime" structs:"UpdateTime"`
}

type DataGameItem struct {
	Sid             int    `json:"Sid" structs:"Sid"`
	Id              string `json:"Id" structs:"Id"`
	Name            string `json:"Name" structs:"Name"`
	UpdateTime      string `json:"UpdateTime" structs:"UpdateTime"`
	GameSid         string `json:"GameSid" structs:"GameSid"`
	ForApi          int    `json:"ForApi" structs:"ForApi"`
	Sort            int    `json:"Sort" structs:"Sort"`
	Enable          bool   `json:"Enable" structs:"Enable"`
	OrderNTD        string `json:"OrderNTD" structs:"OrderNTD"`
	Bonus           string `json:"Bonus" structs:"Bonus"`
	NTD             string `json:"NTD" structs:"NTD"`
	EnableCost      bool   `json:"EnableCost" structs:"EnableCost"`
	Cost            string `json:"Cost" structs:"Cost"`
	AddValueTypeSid string `json:"AddValueTypeSid" structs:"AddValueTypeSid"`
	Note1           string `json:"Note1" structs:"Note1"`
	Note2           string `json:"Note2" structs:"Note2"`
	Count           int    `json:"Count" structs:"Count"`
}

type DataRate struct {
	Sid        int       `json:"Sid" structs:"Sid"`
	Id         string    `json:"Id" structs:"Id"`
	Name       string    `json:"Name" structs:"Name"`
	UpdateTime string    `json:"UpdateTime" structs:"UpdateTime"`
	ListData   CListPair `json:"ListData" structs:"ListData"`
	UserSid    string    `json:"UserSid" structs:"UserSid"`
}

type DataItemCount struct {
	Sid         int    `json:"Sid" structs:"Sid"`
	Id          string `json:"Id" structs:"Id"`
	Name        string `json:"Name" structs:"Name"`
	UpdateTime  string `json:"UpdateTime" structs:"UpdateTime"`
	UserSid     string `json:"UserSid" structs:"UserSid"`
	GameSid     string `json:"GameSid" structs:"GameSid"`
	GameItemSid string `json:"GameItemSid" structs:"GameItemSid"`
	ChangeValue int64  `json:"ChangeValue" structs:"ChangeValue"`
	TotalCount  int64  `json:"TotalCount" structs:"TotalCount"`
	TotalSell   int64  `json:"TotalSell" structs:"TotalSell"`
	OrderSid    string `json:"OrderSid" structs:"OrderSid"`
	GameRate    string `json:"GameRate" structs:"GameRate"`
	Pic0        string `json:"Pic0" structs:"Pic0"`
	Pic1        string `json:"Pic1" structs:"Pic1"`
	Note        string `json:"Note" structs:"Note"`
}

type DataQueryCount struct {
	Sid          int    `json:"Sid" structs:"Sid"`
	Id           string `json:"Id" structs:"Id"`
	GameSid      string `json:"GameSid" structs:"GameSid"`
	GameItemSid  string `json:"GameItemSid" structs:"GameItemSid"`
	Name         string `json:"Name" structs:"Name"`
	CurrentCount int64  `json:"CurrentCount" structs:"CurrentCount"`
	TotalCount   int64  `json:"TotalCount" structs:"TotalCount"`
	TotalSell    int64  `json:"TotalSell" structs:"TotalSell"`
	UpdateTime   string `json:"UpdateTime" structs:"UpdateTime"`
	Note         string `json:"Note" structs:"Note"`
}

type DataPayType struct {
	Sid        int      `json:"Sid" structs:"Sid"`
	Id         string   `json:"Id" structs:"Id"`
	Name       string   `json:"Name" structs:"Name"`
	UpdateTime string   `json:"UpdateTime" structs:"UpdateTime"`
	Value      []string `json:"Value" structs:"Value"`
	SubValue   []string `json:"SubValue" structs:"SubValue"`
	Currency   string   `json:"Currency" structs:"Currency"`
	Sort       int      `json:"Sort" structs:"Sort"`
}

// CListPair is a custom type
type CListPair []struct {
	first  string
	second string
}

//======================

func SetData(data map[string]interface{}, target interface{}) {
	valueOf := reflect.ValueOf(target).Elem()
	typeOf := valueOf.Type()

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		tag := field.Tag.Get("json")
		tags := strings.Split(tag, ",")
		if len(tags) == 0 {
			continue
		}
		val, ok := data[tags[0]]
		if !ok {
			continue
		}

		fieldValue := valueOf.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		fmt.Printf("Field: %s, Value: %v, Type: %T\n", tags[0], val, val) // 添加的调试语句

		switch fieldValue.Interface().(type) {
		case bool:
			boolVal, ok := val.(bool)
			if ok {
				fieldValue.SetBool(boolVal)
			}
		case float64:
			floatVal, ok := val.(float64)
			if ok {
				fieldValue.SetFloat(floatVal)
			}
		case int:
			switch v := val.(type) {
			case string:
				intVal, err := strconv.Atoi(v)
				if err == nil {
					fieldValue.SetInt(int64(intVal))
				}
			case float64:
				fieldValue.SetInt(int64(v))
			}
		case string:
			strVal, ok := val.(string)
			if ok {
				fieldValue.SetString(strVal)
			}
		default:
			fmt.Printf("Unhandled type for field %s: %T\n", tags[0], val) // 添加的调试语句
		}
	}
}
func GetData(target interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	valueOf := reflect.ValueOf(target).Elem()
	typeOf := valueOf.Type()

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		tag := field.Tag.Get("json")
		fieldValue := valueOf.Field(i)

		fmt.Printf("Field: %s, Tag: %s, Value: %v\n", field.Name, tag, fieldValue)

		switch fieldValue.Kind() {
		case reflect.Bool:
			boolVal := fieldValue.Bool()
			intVal := 0
			if boolVal {
				intVal = 1
			}
			data[tag] = intVal
		case reflect.String:
			data[tag] = fieldValue.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			data[tag] = fieldValue.Int()
		case reflect.Float32, reflect.Float64:
			data[tag] = fieldValue.Float()
		default:
			log.Printf("Unhandled type %s for field %s", fieldValue.Kind(), field.Name)
		}
	}

	return data
}

/*
func GetData(target interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	valueOf := reflect.ValueOf(target).Elem()
	typeOf := valueOf.Type()

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		tag := field.Tag.Get("json")
		fieldValue := valueOf.Field(i)

		switch fieldValue.Interface().(type) {
		case bool:
			boolVal := fieldValue.Bool()
			intVal := 0
			if boolVal {
				intVal = 1
			}
			data[tag] = intVal
		case CListPair:
			listVal := fieldValue.Interface().(CListPair)
			data[tag] = listVal.toString()

		default:
			data[tag] = fieldValue.String()
		}
	}

	return data
}
*/

//================

func (clp CListPair) toString() string {
	pairs := make([]string, len(clp))
	for _, pair := range clp {
		pairs = append(pairs, pair.first+":"+pair.second)
	}
	return strings.Join(pairs, ",")
}

func (clp *CListPair) fromString(str string) {
	pairStrs := strings.Split(str, ",")
	for _, pairStr := range pairStrs {
		parts := strings.Split(pairStr, ":")
		if len(parts) == 2 {
			*clp = append(*clp, struct{ first, second string }{parts[0], parts[1]})
		}
	}
}

func (clp CListPair) listFirst() []string {
	keys := make([]string, len(clp))
	for i, pair := range clp {
		keys[i] = pair.first
	}
	return keys
}

func (clp CListPair) listSecond() []string {
	values := make([]string, len(clp))
	for i, pair := range clp {
		values[i] = pair.second
	}
	return values
}

func (clp CListPair) findValue(key string) string {
	for _, pair := range clp {
		if pair.first == key {
			return pair.second
		}
	}
	return ""
}

func ToString(data interface{}) string {
	switch v := data.(type) {
	case string:
		return v
	case bool:
		if v == true {
			return "1"
		} else {
			return "0"
		}
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func main() {
	// Your main program logic can go here
}
