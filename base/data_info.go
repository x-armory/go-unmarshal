package base

import (
	"errors"
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
	BaseTags       *map[int]*FieldTag
}

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
	var baseTags *map[int]*FieldTag
	if tag, e := GetFieldTags(itemStructType, "xm"); e != nil {
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
		baseTags,
	}, nil
}
