package Common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func MapToJsonStr(param map[string]interface{}) string {
	dataType, _ := json.Marshal(param)
	datastring := string(dataType)
	return datastring
}

func JsonStrToMap(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}

	return tempMap
}

// map TO Interface  ,    interface =map
func InterFaceToMap(in interface{}, out *map[string]interface{}) {
	bytes, _ := json.Marshal(in)
	json.Unmarshal(bytes, &out)
}

func InterfaceToString(data interface{}) string {
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

func InterfaceToInt(data interface{}) int {
	switch v := data.(type) {
	case int:
		return v
	case bool:
		if v == true {
			return 1
		} else {
			return 0
		}

	case int64: // 考慮到可能是從 json.Unmarshal 獲得的資料
		return int(v)
	case float64:
		return int(v)
	case string:
		tmp, _ := strconv.Atoi(v)
		return tmp
	default:
		fmt.Printf("invalid type for conversion to int: %T", v)
		return 0
	}
}

type VariantMap map[string]interface{}

func (v *VariantMap) Origin() *map[string]interface{} {
	return (*map[string]interface{})(v)
}

func (v VariantMap) ToInterface() interface{} {
	var re interface{} = v
	return re
}

func InterfaceToStruct(input interface{}, output interface{}) error {

	fmt.Printf("AAA1 %v \n", input)

	data, ok := input.(map[string]interface{})
	if !ok {
		return fmt.Errorf("input is not a map[string]interface{}")
	}

	val := reflect.ValueOf(output)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("output should be a non-nil pointer")
	}
	val = val.Elem()

	if val.Kind() != reflect.Struct {
		fmt.Printf("Expected a struct but got: %v\n", val.Kind())
		return fmt.Errorf("output should be a pointer to a struct")
	}

	missingFields := []string{} // 创建一个切片来记录不存在的字段

	for k, v := range data {
		if v == nil { // 过滤nil值
			continue
		}
		field, exists := val.Type().FieldByName(k) // 检查字段是否存在
		if !exists {
			missingFields = append(missingFields, k) // 如果字段不存在，将其添加到切片中
			continue
		}

		if field.Type.Kind() == reflect.Int && reflect.TypeOf(v).Kind() == reflect.String {
			// 如果结构体字段是int类型，但数据中的值是string类型
			intValue, err := strconv.Atoi(v.(string))
			if err != nil {
				return fmt.Errorf("error converting field %s from string to int: %v", k, err)
			}
			data[k] = intValue
		}

		if field.Type.Kind() == reflect.String && reflect.TypeOf(v).Kind() == reflect.Int64 {
			// Convert int64 to string for fields like GameSid
			data[k] = strconv.FormatInt(v.(int64), 10)
		} else if field.Type.Kind() == reflect.Bool && reflect.TypeOf(v).Kind() == reflect.Int64 {
			// Convert non-zero int64 to true and zero to false for fields like Enable and EnableCost
			if v.(int64) == 0 {
				data[k] = false
			} else {
				data[k] = true
			}
		}

	}

	if len(missingFields) > 0 {
		//return fmt.Errorf("the following fields do not exist in the struct: %v", missingFields)
		fmt.Printf("the following fields do not exist in the struct: %v\n", missingFields)

	}
	fmt.Printf("AAA2 %v \n", data)

	return MapToStruct(data, output)
}

func MapToStructByName(data map[string]interface{}, target interface{}) {
	val := reflect.ValueOf(target).Elem()

	for key, value := range data {
		field := val.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		}
	}
}

func ConvertSid(data map[string]interface{}, outputStruct interface{}) error {
	// 首先，检查data["Sid"]的类型
	var sidInt64 int64
	sid, exists := data["Sid"]

	if !exists {
		return nil
	}

	switch sidValue := sid.(type) {
	case string:
		var err error
		sidInt64, err = strconv.ParseInt(sidValue, 10, 64)
		if err != nil {
			return nil
		}
		data["Sid"] = sidInt64
	case int:
		sidInt64 = int64(sidValue)
		data["Sid"] = sidInt64
	case int64:
		sidInt64 = sidValue
	default:
		return nil
	}

	// 然后，检查outputStruct中Sid字段的类型，并进行适当的转换
	val := reflect.ValueOf(outputStruct).Elem()
	sidField := val.FieldByName("Sid")

	if sidField.IsValid() {
		switch sidField.Kind() {
		case reflect.String:
			sidField.SetString(strconv.FormatInt(sidInt64, 10))
		case reflect.Int, reflect.Int64:
			sidField.SetInt(sidInt64)
		default:
			return nil
		}
	}

	return nil
}

func ConvertFloatToInt(data map[string]interface{}) {
	// for key, value := range data {
	// 	switch v := value.(type) {
	// 	case float64:
	// 		// 確保轉換不會丟失資料
	// 		data[key] = int64(v)
	// 	}
	// }
}

func MapToStruct(data map[string]interface{}, outputStruct interface{}) error {

	//ConvertFloatToInt(data)

	ConvertSid(data, outputStruct)

	// fmt.Println("AAAAAAA1");
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, outputStruct)
	if err != nil {
		return err
	}

	return nil
}

func StructToMap(item interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(item)

	// 如果傳入的是指針，獲取其元素值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 確保傳入的是結構體
	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		// 獲取 struct tag `json` 的值作為key
		key := typ.Field(i).Tag.Get("json")
		if key == "" {
			continue
		}

		result[key] = val.Field(i).Interface()
	}
	return result
}

func StringToInt64(str string) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		// 這裡只是一個例子，你可能想要有不同的錯誤處理方式。
		fmt.Println("Error parsing string:", err)
		return 0
	}
	return val
}

func Int64ToString(val int64) string {
	return strconv.FormatInt(val, 10)
}
func StringToFloat64(str string) float64 {
	value, err := strconv.ParseFloat(str, 64) // 第二個參數64表示轉換成float64
	if err != nil {
		fmt.Println("轉換錯誤:", err)
	}
	return value
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64) // 'f'表示普通的浮點格式。-1表示將所有小數位都列出。64表示f是float64。
}

func StringToList(str string, separator string, length ...int) []string {
	list := strings.Split(str, separator)
	if len(length) > 0 {
		for i := len(list); i < length[0]; i++ {
			list = append(list, "")
		}
	}

	return list
}

func ListToString(list []string, separator string) string {
	return strings.Join(list, separator)
}

func ToListMap(listData []interface{}) ([]map[string]interface{}, error) {
	var convertedListData []map[string]interface{}

	for _, item := range listData {
		if mapData, ok := item.(map[string]interface{}); ok {
			convertedListData = append(convertedListData, mapData)
		} else {
			return nil, fmt.Errorf("Unexpected data format in ListData")
		}
	}

	return convertedListData, nil
}

// -------
func (v VariantMap) Map(key string) VariantMap {
	if entry, ok := v[key].(map[string]interface{}); ok {
		return entry
	}
	return nil
}

func (v VariantMap) String(key string) string {
	if value, found := v[key]; found {
		switch typedValue := value.(type) {
		case string:
			return typedValue
		case int:
			return strconv.Itoa(typedValue)
		case float64: // Go 通常解析JSON為float64，但如果需要支援其他浮點型態，您可以增加更多的cases
			return strconv.FormatFloat(typedValue, 'f', -1, 64)
		case bool:
			if typedValue {
				return "1"
			}
			return "0"
		}
	}
	return ""
}

func (v VariantMap) ToList(key string) []map[string]interface{} {
	if value, found := v[key]; found {
		if list, isList := value.([]map[string]interface{}); isList {
			return list
		}
	}
	return nil
}

func TimeUtc8Str() string {
	//re := time.Now().UTC().Add(time.Hour * 8).Format("20060102150405")

	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return ""
	}

	// 獲取當前台灣時區的時間
	taiwanTime := time.Now().In(location)

	// 格式化時間為字符串
	re := taiwanTime.Format("20060102150405")

	return re
}

func DateUtc8Str(addDays int) string {
	// 加載亞洲/台北時區
	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return ""
	}

	// 獲取當前台灣時區的時間
	taiwanTime := time.Now().In(location)

	// 添加指定天數
	if addDays != 0 {
		taiwanTime = taiwanTime.Add(time.Hour * 24 * time.Duration(addDays))
	}

	// 格式化時間為字符串
	re := taiwanTime.Format("20060102")

	return re
}

func PrintHHMMSS(st string) string {
	re := time.Now().UTC().Add(time.Hour * 8).Format("15:04:05.999")
	fmt.Printf("[%s] %s \n", re, st)
	return re
}

func PrintMap(sTag string, d map[string]interface{}) {

	for k, v := range d {
		fmt.Printf(sTag+" Key: %s, Value: %v, Type: %T\n", k, v, v)
	}

}

func PrintStruct(prefix string, s interface{}) {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Struct {
		fmt.Println("The provided item is not a struct!")
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		fmt.Printf("%s Field Name: %s, Field Value: %v, Type: %T\n", prefix, field.Name, value.Interface(), value.Interface())
	}
}

func PrintInterface(prefix string, data interface{}) {
	val := reflect.ValueOf(data)
	typ := val.Type()

	switch typ.Kind() {
	case reflect.Map:
		// Iterate over map keys and values
		for _, key := range val.MapKeys() {
			value := val.MapIndex(key)
			fmt.Printf("%s Key: %v, Value: %v, Type: %T\n", prefix, key.Interface(), value.Interface(), value.Interface())
		}
	case reflect.Struct:
		// Iterate over struct fields
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			value := val.Field(i)
			fmt.Printf("%s Field Name: %s, Field Value: %v, Type: %T\n", prefix, field.Name, value.Interface(), value.Interface())
		}
	default:
		fmt.Printf("%s Value: %v, Type: %T\n", prefix, data, data)
	}
}
