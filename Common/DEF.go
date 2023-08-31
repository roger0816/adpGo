package Common

const (
	SET_VALUE  = 9901
	GET_VALUE  = 9931
	QUERY_INFO = 10
	LOGIN      = 1000

	ADD_USER   = 1001
	DEL_USER   = 1002
	EDIT_USER  = 1003
	QUERY_USER = 1031

	ADD_CUSTOMER     = 1101
	DEL_CUSTOMER     = 1102
	EDIT_CUSTOMER    = 1103
	QUERY_CUSTOMER   = 1131

	REPLACE_GAME_INFO        = 1201
	DEL_GAME_INFO            = 1202
	QUERY_CUSTOMER_GAME_INFO = 1231

	ADD_CUSTOMER_COST     = 1301
	QUERY_CUSTOMER_COST   = 1331
	LAST_CUSTOMER_COST_ID = 1332

	ADD_GAME_LIST   = 1401
	EDIT_GAME_LIST  = 1402
	DEL_GAME_LIST   = 1403
	QUERY_GAME_LIST = 1431
	QUERY_GAME_RATE = 1432

	ADD_GAME_ITEM   = 1501
	EDIT_GAME_ITEM  = 1502
	DEL_GAME_ITEM   = 1503
	QUERY_GAME_ITEM = 1531

	ADD_BULLETIN   = 1601
	EDIT_BULLETIN  = 1602
	DEL_BULLETIN   = 1603
	QUERY_BULLETIN = 1631

	ADD_CUSTOM_CLASS  = 1701
	EDIT_CUSTOM_CLASS = 1702
	DEL_CUSTOM_CLASS  = 1703
	// CUSTOM_DEBIT 待移除, 改用DEBIT_CLASS
	//ADD_CUSTOM_DEBIT  = 1711
	//EDIT_CUSTOM_DEBIT = 1712
	//DEL_CUSTOM_DEBIT  = 1713

	QUERY_CUSTOM_CLASS = 1731
	QUERY_CUSTOM_DEBIT = 1732

	ADD_FACTORY_CLASS   = 1801
	EDIT_FACTORY_CLASS  = 1802
	DEL_FACTORY_CLASS   = 1803
	QUERY_FACTORY_CLASS = 1831

	ADD_PAY_TYPE   = 1901
	DEL_PAY_TYPE   = 1902
	EDIT_PAY_TYPE  = 1903
	QUERY_PAY_TYPE = 1931

	//ADD_GROUP   = 2001
	//DEL_GROUP   = 2002
	//EDIT_GROUP  = 2003
	//QUERY_GROUP = 2031

	REPLACE_ORDER   = 2101
	QUERY_ORDER     = 2131
	LAST_ORDER_ID   = 2132
	LAST_ORDER_NAME = 2133

	ADD_BOUNS   = 2201
	QUERY_BOUNS = 2231

	ADD_SCHEDULE   = 2301
	QUERY_SCHEDULE = 2331

	ADD_EXCHANGE   = 2401
	QUERY_EXCHANGE = 2431

	ADD_PRIMERATE   = 2501
	QUERY_PRIMERATE = 2531

	UPLOAD_PIC = 2601
	QUERY_PIC  = 2631

	ADD_ITEM_COUNT   = 2701
	EDIT_ITEM_COUNT  = 2702
	DEL_ITEM_COUNT   = 2703
	QUERY_ITEM_COUNT = 2731

	ADD_DEBIT_CLASS   = 2811
	EDIT_DEBIT_CLASS  = 2812
	DEL_DEBIT_CLASS   = 2813
	QUERY_DEBIT_CLASS = 2831

	PAY_ADD_COST = 3002 

	PAY_ORDER = 3003

	QUERY_MIX = 6031

	API_REQUSET = 7000
)

type SqlTable struct{}

var SQL_TABLE = SqlTable{}

func (s SqlTable) ExchangeRate() string     { return "ExchangeRate" }
func (s SqlTable) PrimeCostRate() string    { return "PrimeCostRate" }
func (s SqlTable) GameList() string         { return "GameList" }
func (s SqlTable) GameRate() string         { return "GameRate" }
func (s SqlTable) GameItem() string         { return "GameItem" }
func (s SqlTable) Bulletin() string         { return "Bulletin" }
func (s SqlTable) CustomerMoney() string    { return "CustomerMoney" }
func (s SqlTable) CustomerClass() string    { return "CustomerClass" }
func (s SqlTable) DebitClass() string       { return "DebitClass" }
func (s SqlTable) GroupData() string        { return "GroupData" }
func (s SqlTable) FactoryClass() string     { return "FactoryClass" }
func (s SqlTable) UserData() string         { return "UserData" }
func (s SqlTable) CustomerData() string     { return "CustomerData" }
func (s SqlTable) CustomerCost() string     { return "CustomerCost" }
func (s SqlTable) CustomerGameInfo() string { return "CustomerGameInfo" }
func (s SqlTable) OrderData() string        { return "OrderData" }
func (s SqlTable) PayType() string          { return "PayType" }
func (s SqlTable) UserBonus() string        { return "UserBonus" }
func (s SqlTable) Schedule() string         { return "Schedule" }
func (s SqlTable) PicData() string          { return "PicData" }
func (s SqlTable) GameItemCount() string    { return "GameItemCount" }
func (s SqlTable) QueryCount() string       { return "QueryCount" }
func (s SqlTable) Settings() string         { return "Settings" }
