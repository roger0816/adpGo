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
		fmt.Errorf("invalid type for conversion to int: %T", v)
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

func (v VariantMap) ToStruct(target interface{}) {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		fmt.Println("Invalid target type")
		return
	}

	for fieldName, fieldValue := range v {
		field := targetValue.Elem().FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(fieldValue))
		}
	}
}

func mapToStruct(data map[string]interface{}, target interface{}) {
	val := reflect.ValueOf(target).Elem()

	for key, value := range data {
		field := val.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		}
	}
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
