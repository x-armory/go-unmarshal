package base

import (
	"errors"
	"fmt"
	"reflect"
)

// 负责处理原始数据ItemAllFieldsValue，转化为目标对象
type DataLoader struct {
	Data           interface{}
	Target         reflect.Value
	DataType       reflect.Type
	DataIsSlice    bool
	ItemType       reflect.Type
	ItemIsPtr      bool
	ItemStructType reflect.Type
	BaseTags       *map[int]*FieldTag

	// FetchDataFunc返回chan，每个数据表示一个item数据
	writeData bool
	itemCache reflect.Value
	filters   []ItemFilter
}

// 每个map表示一个item数据，key表示field index，value表示内容字符串，
// 由FieldUnmarshalTag负责解析value并生成元素对象，
// key包含-1表示迭代结束。
type ItemAllFieldsValue *map[int]string
type ItemValueChan chan ItemAllFieldsValue

// 处理每个元素，如果任意一个返回false，中断处理
type ItemFilter func(item interface{}) bool

// 创建DataLoader。
// param data，数据存储目标对象指针；
// param itemValueChan；
// param writeData，为true则将每个元素放入cache并写入data指针，否则不缓存也不写入data；
// param filters，处理每个元素，如果任一返回false，中断处理，不论data元素是不是地址类型，filter入参总是*Model指针；
// Tips 将cacheAllItems设为false，并通过filter处理每个元素，适合用来处理超大列表
func NewDataLoader(tagName string, data interface{}, writeData bool, itemFilter ...ItemFilter) (*DataLoader, error) {
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
	if tag, e := GetFieldTags(itemStructType, tagName); e != nil {
		return nil, e
	} else {
		baseTags = tag
	}
	if len(*baseTags) == 0 {
		return nil, errors.New(fmt.Sprintf("%s tag not found", tagName))
	}

	// 元素缓存，不论目标是单个对象还是数组，统一用数组缓存，方便统一处理
	var cache = reflect.MakeSlice(reflect.SliceOf(itemType), 0, 0)

	return &DataLoader{
		Data:           data,
		Target:         target,
		DataType:       dataType,
		DataIsSlice:    dataIsSlice,
		ItemType:       itemType,
		ItemIsPtr:      itemIsPtr,
		ItemStructType: itemStructType,
		BaseTags:       baseTags,
		writeData:      writeData,
		itemCache:      cache,
		filters:        itemFilter,
	}, nil
}

func (i *DataLoader) LoadData(itemValueChan ItemValueChan) error {
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
				i.itemCache = reflect.Append(i.itemCache, *itemPtr)
			} else {
				i.itemCache = reflect.Append(i.itemCache, reflect.Indirect(*itemPtr))
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
		i.Target.Set(i.itemCache)
	} else if i.itemCache.Len() > 0 {
		i.Target.Set(i.itemCache.Index(0))
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
