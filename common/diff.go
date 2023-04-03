package common

import "reflect"

type Diff interface {
	GetKey() string
	Values() []string
}

//DiffContrast 泛型语法 更方便  需要go 1.18支持
//第一个返回切片表示对比数据新增的数据
//第二个返回切片表示对比数据已修改的数据
//第三个返回切片表示对比数据删除的数据
func DiffContrast[T Diff](oldList, newList []T) ([]T, []T, []T) {
	if len(oldList) == 0 { //旧的数据没有 那全是新增的数据
		return newList, nil, nil
	}
	oldMap := Map[string, []string]{}
	for _, value := range oldList {
		oldMap[value.GetKey()] = value.Values()
	}
	newMap := Map[string, []string]{}
	for _, value := range newList {
		newMap[value.GetKey()] = value.Values()
	}

	addList := make([]T, 0)
	updateList := make([]T, 0)
	deleteList := make([]T, 0)

	//遍历新的
	for _, value := range newList {
		//旧的数据中找到了  如果value不一样 那就是修改的数据
		if old, ok := oldMap[value.GetKey()]; ok {
			for i, val := range old {
				if val != value.Values()[i] {
					updateList = append(updateList, value)
					break
				}
			}

		} else { //旧的中找不到 那就是新增的数据
			addList = append(addList, value)
		}
	}

	//遍历旧的
	for _, value := range oldList {
		//新的数据中找不到 那就是删除后的数据
		if _, ok := newMap[value.GetKey()]; !ok {
			deleteList = append(deleteList, value)
		}
	}

	return addList, updateList, deleteList
}

// DiffField
// 对比结构体中字段的差异。仅支持string,uint,uint8,uint16,uint32,uint64,int,int8,int16,int32,int64,bool,float32,float64
func DiffField[T any](oldData, newData *T) *T {
	oldDataMap := Map[int, any]{}

	oldValueOf := reflect.ValueOf(oldData).Elem()
	for i := 0; i < oldValueOf.NumField(); i++ {
		value := oldValueOf.Field(i)
		if value.CanInterface() {
			oldDataMap[i] = value.Interface()
		}
	}

	resData := new(T)
	resValue := reflect.ValueOf(resData).Elem()
	newValueOf := reflect.ValueOf(newData).Elem()
	for i := 0; i < newValueOf.NumField(); i++ {
		field := newValueOf.Field(i)
		if !field.CanInterface() {
			continue
		}
		value := field.Interface()
		if oldValue, ok := oldDataMap[i]; ok {
			resField := resValue.Field(i)
			if value != oldValue && resField.CanSet() {
				switch resField.Kind() {
				case reflect.String:
					resField.SetString(field.String())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					resField.SetUint(field.Uint())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					resField.SetInt(field.Int())
				case reflect.Bool:
					resField.SetBool(field.Bool())
				case reflect.Float64, reflect.Float32:
					resField.SetFloat(field.Float())
				default:
					continue
				}
			}
		}
	}
	return resData
}

// DiffMap 对比两个map的key,value是否都一样，是一样返回tr
func DiffMap[T, V comparable](oldMap, newMap map[T]V) bool {
	for oldKey, oldVal := range oldMap {
		if newVal, ok := newMap[oldKey]; !ok {
			return false
		} else {
			if newVal != oldVal {
				return false
			}
		}
	}

	for newKey, newValue := range newMap {
		if oldVal, ok := oldMap[newKey]; !ok {
			return false
		} else {
			if newValue != oldVal {
				return false
			}
		}
	}

	return true
}
