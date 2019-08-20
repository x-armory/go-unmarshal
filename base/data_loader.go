package base

import (
	"errors"
	"fmt"
	"reflect"
)

// 负责处理原始数据ItemAllFieldsValue，转化为目标对象
// writeData，为true则将每个元素放入cache并写入data指针，否则不缓存也不写入data；
// filters，处理每个元素，如果任一返回false，中断处理，不论data元素是不是地址类型，filter入参总是*Model指针；
// Tips: 将writeData设为false，通过filter处理每个元素，适合用来处理超大列表
type DataLoader struct {
	*DataInfo
	// FetchDataFunc返回chan，每个数据表示一个item数据
	writeData bool
	filters   []ItemFilter
}

func NewDataLoader(v interface{}, writeData bool, filters ...ItemFilter) (*DataLoader, error) {
	info, e := GetDataInfo(v)
	if e != nil {
		return nil, e
	}
	return &DataLoader{
		info,
		writeData,
		filters,
	}, nil
}

// 每个map表示一个item数据，key表示field index，value表示内容字符串，
// 由FieldUnmarshalTag负责解析value并生成元素对象，
// key包含-1表示迭代结束。
type ItemAllFieldsValue *map[int]string
type ItemValueChan chan ItemAllFieldsValue

// 处理每个元素，如果任意一个返回false，中断处理
type ItemFilter func(item interface{}) bool

func (i *DataLoader) LoadData(itemValueChan ItemValueChan) error {
	// 元素缓存，不论目标是单个对象还是数组，统一用数组缓存，方便统一处理
	var cache reflect.Value
	if i.writeData {
		cache = reflect.MakeSlice(reflect.SliceOf(i.ItemType), 0, 0)
	}
loadLoop:
	for itemData := range itemValueChan {
		if _, ok := (*itemData)[-1]; ok {
			break loadLoop
		}
		itemPtr, err := i.buildItemValue(itemData)
		if err != nil {
			return err
		}
		if itemPtr == nil {
			break loadLoop
		}
		if i.writeData {
			if i.ItemIsPtr {
				cache = reflect.Append(cache, *itemPtr)
			} else {
				cache = reflect.Append(cache, reflect.Indirect(*itemPtr))
			}
		}
		for _, f := range i.filters {
			if !f(itemPtr.Interface()) {
				break loadLoop
			}
		}
		if !i.DataIsSlice {
			break loadLoop

		}
	}

	if !i.writeData {
		return nil
	}
	if i.DataIsSlice {
		i.Target.Set(cache)
	} else if cache.Len() > 0 {
		i.Target.Set(cache.Index(0))
	}
	return nil
}

// 根据获取的数据组装对象元素
func (i *DataLoader) buildItemValue(itemData ItemAllFieldsValue) (*reflect.Value, error) {
	if itemData == nil {
		return nil, nil
	}
	value := reflect.New(i.ItemStructType)
	valueTarget := reflect.Indirect(value)
	hasFieldValue := false
	for fieldIndex, fieldValueString := range *itemData {
		fieldTag, ok := (*i.BaseTags)[fieldIndex]
		if !ok {
			return nil, errors.New(fmt.Sprintf("bad field index %d for type %s", fieldIndex, i.ItemStructType.Name()))
		}
		if v, err := fieldTag.Parse(fieldValueString); err != nil {
			return nil, err
		} else if v.IsValid() {
			valueTarget.Field(fieldIndex).Set(v)
			hasFieldValue = true
		}
	}
	if hasFieldValue {
		return &value, nil
	} else {
		return nil, nil
	}
}
