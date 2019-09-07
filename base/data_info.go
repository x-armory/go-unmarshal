package base

import (
	"errors"
	"fmt"
	"reflect"
)

type DataInfo struct {
	Data           interface{}
	Target         reflect.Value
	DataType       reflect.Type
	DataIsSlice    bool
	ItemType       reflect.Type
	ItemIsPtr      bool
	ItemStructType reflect.Type
	BaseTags       FieldTagMap
}

// data必须可写，包括指针、数组，必须是Struct 或者 Slice of Struct
func GetDataInfo(data interface{}) (*DataInfo, error) {
	// 校验非空
	if data == nil {
		return nil, errors.New("data can not be nil")
	}
	// 校验data类型
	dataValue := reflect.ValueOf(data)
	target := reflect.Indirect(dataValue)
	if !target.CanSet() || target.Kind() != reflect.Slice && target.Kind() != reflect.Struct {
		return nil, errors.New("needs a pointer to a slice or a struct")
	}

	// 获取data、data元素、data元素struct类型
	dataType := reflect.TypeOf(data).Elem()
	dataIsSlice := dataType.Kind() == reflect.Slice
	itemType := dataType
	if dataIsSlice {
		itemType = dataType.Elem()
	}
	itemIsPtr := itemType.Kind() == reflect.Ptr
	itemStructType := itemType
	if itemIsPtr {
		itemStructType = itemType.Elem()
	}
	if itemStructType.Kind() != reflect.Struct {
		return nil, errors.New("needs a pointer to a slice or a struct")
	}

	// 解析base tags
	var baseTags *FieldTagMap
	if tag, e := GetFieldTags(itemStructType, "xm", nil); e != nil {
		return nil, e
	} else {
		baseTags = tag
	}

	return &DataInfo{
		data,
		target,
		dataType,
		dataIsSlice,
		itemType,
		itemIsPtr,
		itemStructType,
		*baseTags,
	}, nil
}

// 根据 map[fieldIndex]fieldStringValue 组装对象元素
func BuildItemValue(dataInfo *DataInfo, itemData *FieldValueMap) (*reflect.Value, error) {
	if itemData == nil {
		return nil, nil
	}
	value := reflect.New(dataInfo.ItemStructType)
	valueTarget := reflect.Indirect(value)
	hasFieldValue := false
	for fieldIndex, fieldValue := range *itemData {
		_, ok := dataInfo.BaseTags[fieldIndex]
		if !ok {
			return nil, errors.New(fmt.Sprintf("bad field index %d for type %s", fieldIndex, dataInfo.ItemStructType.Name()))
		}
		valueTarget.Field(fieldIndex).Set(fieldValue)
		hasFieldValue = true
	}
	if hasFieldValue {
		return &value, nil
	} else {
		return nil, nil
	}
}
