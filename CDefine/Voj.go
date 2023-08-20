package CDefine

import (
	"strings"
)

type DataObj struct {
	Sid        string
	Id         string
	Name       string
	UpdateTime string
}

type UserData struct {
	DataObj
	Password   string
	Cid        string
	Lv         int
	ParentId   string
	StartDay   string
	BirthDay   string
	Tel        string
	Email      string
	Note1      string
	Note2      string
	Note3      string
	CreateTime string
}

type DataUserBonus struct {
	Sid      string
	UserSid  string
	OrderSid string
	Bonus    string
	Pay      string
	AddBonus string
	AddPay   string
	UpdateTime string
}

type CustomerData struct {
	DataObj
	Vip       string
	Class     string
	Money     string
	Line      string
	Currency  string
	PayType   string
	PayInfo   string
	UserSid   string
	Num5      string
	Note1     string
	Note2     string
}

type CustomerMoney struct {
	DataObj
	Money    string
	Currency string
}

type CustomerCost struct {
	Sid           string
	CustomerSid   string
	OrderId       string
	Currency      string
	ChangeValue   string
	OriginCurrency string
	OriginValue   string
	DebitSid      string
	DebitNote     string
	Rate          string
	AddRate       string
	IsAddCost     bool
	Total         string
	UserSid       string
	OrderTime     string
	UpdateTime    string
	Pic0          string
	Pic1          string
	Note0         string
	Note1         string
}

type DataCustomerClass struct {
	DataObj
	Sort int
	Type string
	Note1 string
	Note2 string
}

type DebitClass struct {
	DataObj
	Sort     int
	Currency string
	Note1    string
	Note2    string
}

type GroupData struct {
	DataObj
	Type      string
	Value     string
	Blob      []byte
	Note1     string
	Note2     string
	NoteBlob  []byte
}

type CustomerGameInfo struct {
	DataObj
	CustomerSid    string
	CustomerId     string
	GameSid        string
	LoginType      string
	LoginAccount   string
	LoginPassword  string
	ServerName     string
	Characters     string
	LastTime       string
	Note1          string
}

type OrderData struct {
	DataObj
	CustomerSid     string
	UiRecord        string
	Step            string
	StepTime        []string
	User            []string
	Owner           string
	PaddingUser     string
	GameSid         string
	Item            string
	Cost            string
	Bouns           string
	PayType         string
	CanSelectPayType string
	GameRate        string
	ExRateSid       string
	PrimeRateSid    string
	Money           []string
	Note0           []string
	ListCost        []string
	ListBouns       []string
	ItemInfo        []string
	Note1           string
	Note2           string
	Note3           string
	Note4           string
	Note5           string
	Pic0            string
	Pic1            string
	OrderDate       string
	OrderTime       string
	CustomerName    string
	Currency        string
}

type DataFactory struct {
	DataObj
	Currency     string
	PayTypdSid   []string
}

type DataGameList struct {
	DataObj
	Enable    bool
	GameRate  float64
	UserSid   string
	SellNote  string
}

type DataGameRate struct {
	Sid       string
	GameSid   string
	GameName  string
	Rate      string
	UserSid   string
	UpdateTime string
}

type DataGameItem struct {
	DataObj
	GameSid      string
	Sort         int
	Enable       bool
	OrderNTD     string
	Bonus        string
	NTD          string
	EnableCost   bool
	Cost         string
	AddValueTypeSid string
	Note1        string
	Note2        string
	Count        int
}

type DataRate struct {
	DataObj
	listData CListPair
	UserSid  string
}

type DataItemCount struct {
	DataObj
	UserSid      string
	GameSid      string
	GameItemSid  string
	ChangeValue  int64
	TotalCount   int64
	TotalSell    int64
	OrderSid     string
	GameRate     string
	Pic0         string
	Pic1         string
	Note         string
}

type DataPayType struct {
	DataObj
	Value     []string
	SubValue  []string
	Currency  string
	Sort      int
}

// CListPair is a custom type
type CListPair []struct {
	first  string
	second string
}

func (clp CListPair) toString() string {
	pairs := make([]string, len(clp))
	for i, pair := range clp {
		pairs[i] = pair.first + ":" + pair.second
	}
	return strings.Join(pairs, ",")
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

func main() {
	// Your main program logic can go here
}
