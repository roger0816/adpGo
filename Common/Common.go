package Common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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

func MapToStruct(data map[string]interface{}, outputStruct interface{}) error {

	// 對於特定的鍵（例如Sid），進行預處理轉換
	// if sidValue, ok := data["Sid"]; ok {
	// 	switch sid := sidValue.(type) {
	// 	case int, int64:
	// 		data["Sid"] = fmt.Sprintf("%d", sid)
	// 		// case float64:
	// 		//     data["Sid"] = fmt.Sprintf("%.0f", sid)
	// 	case float64:
	// 		data["Sid"] = fmt.Sprintf("%d", int64(v)) // 轉為 int64，移除小數部分
	// 	}
	// }

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
	re := time.Now().UTC().Add(time.Hour * 8).Format("20060102150405")
	return re
}

func PrintHHMMSS(st string) string {
	re := time.Now().UTC().Add(time.Hour * 8).Format("15:04:05.999")
	fmt.Printf("[%s] %s \n", re, st)
	return re
}
