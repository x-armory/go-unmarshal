package base

import (
	"errors"
	"fmt"
	"reflect"
)

// 负责遍历所有变量，反序列化数据，调用ItemFilters处理每个元素，并最终根据writeData开关决定是否写入DataInfo.Target；
// 将writeData设为false，通过filter处理每个元素，可避免内存浪费，适合用来处理超大列表；
// WriteData，是否将数据写入；DataInfo.Target，关闭后可以减少内存占用；
// DataInfo，写入目标对象信息，包含对象的类型、指针、字段注释等；
// VarOrder，用于控制多个变量嵌套循环顺序，如果数组为空，按随机顺序嵌套循环；
// FieldTagReadValueFuncMap，结构为map[tag]FieldTagReadValueFunc，定义每种tag到读取方法，tag包括excel、xpath等；
// ItemFilters，处理每个的原始数据、转换对象，如果任一返回false，中断处理，不论data元素是不是地址类型
// filter入参总是*interface{}指针，在校验原始数据时，入仓类型是*FieldValueMap，在校验元素对象时，类型是*YourModel，非已知类型最好返回(Forward,0)；
type DataLoader struct {
	Data            interface{}
	WriteData       bool
	VarOrder        []string
	ExitNoDataTimes int
	ReadValueFunc   map[string]FieldTagReadValueFunc
	ItemFilters
}

// 迭代DataLoader.DataInfo.BaseTags.MergeVars
// 根据DataLoader.DataInfo、Vars、文档读取方法，获取对象实例
// 调用DataLoader.filters，处理元素，并控制读取流程
func (loader *DataLoader) Load() error {
	// default values
	if loader.ExitNoDataTimes < 1 {
		loader.ExitNoDataTimes = 1
	} else if loader.ExitNoDataTimes > 10 {
		loader.ExitNoDataTimes = 10
	}
	// 校验参数
	if loader.Data == nil {
		return errors.New("data is nil")
	}
	if len(loader.ReadValueFunc) == 0 {
		return errors.New("loader.ReadValueFunc is empty")
	}
	if !loader.WriteData && len(loader.ItemFilters) == 0 {
		return errors.New("need at least one ItemFilter while WriteData=false")
	}
	// gen data info
	dataInfo, e := GetDataInfo(loader.Data)
	if e != nil {
		return e
	}
	// 支持多tag，不再筛选tag，比如zip+xls
	//tags := loader.BaseTags.Filter(func(tag *FieldTag) bool {
	//	return tag.Schema == loader.tag
	//})
	tags := &dataInfo.BaseTags
	vars := tags.MergeVars().List(loader.VarOrder...)
	// validate tag and loader.ReadValueFunc
	for _, tag := range *tags {
		if _, ok := loader.ReadValueFunc[tag.Schema]; !ok {
			return errors.New("no readValueFunc specified for " + tag.Schema + " schema")
		}
	}

	// 创建缓存
	var cache reflect.Value
	if loader.WriteData {
		cache = reflect.MakeSlice(reflect.SliceOf(dataInfo.ItemType), 0, 0)
	}

	var noMoreValueTimes = 0
	// 遍历所有可能的变量，读取内容，过滤检查，设置缓存
	for vars.Reset(); vars.IsValid(); {
		// 设置变量
		tags.SetValues(vars)
		// 判断带变量到tag是否匹配读取到值，只要一个字段匹配到就算有值
		var isVarMatchedValue = false
		// 读取所有字段内容
		var allFieldValue = make(FieldValueMap)
		// 调用场景定制ReadItemFieldFunc读取所有字段内容
		for fieldIdx, tag := range *tags {
			// check path filled with vars
			if tag.PathFilled == "" {
				println("[WARN] " + tag.FieldName + " pathFilled is empty, " + tag.Path + ", " + vars.ToString())
				continue
			}
			// read field value
			if s, err := loader.ReadValueFunc[tag.Schema](tag, vars); err != nil {
				println("[WARN]", err.Error())
			} else {
				value, err := tag.Parse(s)
				if err != nil {
					return err
				}
				if value.IsValid() {
					allFieldValue[fieldIdx] = value
					if len(*tag.Vars) > 0 {
						// 带变量的tag匹配到值 清零计数
						isVarMatchedValue = true
						noMoreValueTimes = 0
					}
				}
			}
		}
		// 对所有原始数据进行过滤、校验
		flow, deep := loader.ItemFilters.Filter(&allFieldValue, vars)
		if flow == Break {
			vars.Break(deep)
			continue
		} else if flow == Continue {
			vars.Continue(deep)
			continue
		}
		// 转换为对象实例
		item, e := BuildItemValue(dataInfo, &allFieldValue)
		if e != nil {
			return e
		}
		if item == nil {
			println("[WARN] item is nil, ", vars.ToString())
		}
		// 增加计数
		if !isVarMatchedValue || item == nil {
			noMoreValueTimes++
		}
		// 带变量的tag是否达到退出次数设置
		//  if noMoreValueTimes == loader.ExitNoDataTimes exit
		//  if noMoreValueTimes>0 && noMoreValueTimes == loader.ExitNoDataTimes-1 break(n)
		//  if noMoreValueTimes>0 && noMoreValueTimes < loader.ExitNoDataTimes-1 continue(n)
		if noMoreValueTimes > 0 {
			if noMoreValueTimes < loader.ExitNoDataTimes-1 {
				deep := len(*vars) - 1
				println(fmt.Sprintf("[WARN] no data %d times, loader.ExitNoDataTimes=%d, continue(deep=%d), %s", noMoreValueTimes, loader.ExitNoDataTimes, deep, vars.ToString()))
				vars.Continue(deep)
				continue
			} else if noMoreValueTimes == loader.ExitNoDataTimes-1 {
				deep := len(*vars) - 1
				println(fmt.Sprintf("[WARN] no data %d times, loader.ExitNoDataTimes=%d, break(deep=%d), %s", noMoreValueTimes, loader.ExitNoDataTimes, deep, vars.ToString()))
				vars.Break(deep)
				continue
			} else {
				println(fmt.Sprintf("[WARN] no data %d times, loader.ExitNoDataTimes=%d, break(deep=%d), %s", noMoreValueTimes, loader.ExitNoDataTimes, 0, vars.ToString()))
				vars.Break(0)
				continue
			}
		}
		// 对转换结果进行过滤、校验
		var itemInterface interface{}
		if item == nil {
			itemInterface = nil
		} else {
			itemInterface = item.Interface()
		}
		flow, deep = loader.ItemFilters.Filter(itemInterface, vars)
		if flow == Break {
			vars.Break(deep)
			continue
		} else if flow == Continue {
			vars.Continue(deep)
			continue
		}
		// 通过filter校验，写入缓存
		if loader.WriteData {
			if dataInfo.ItemIsPtr {
				cache = reflect.Append(cache, *item)
			} else {
				cache = reflect.Append(cache, reflect.Indirect(*item))
			}
		}
		vars.Next()
	}
	// 如果写数据关闭，直接退出
	if !loader.WriteData {
		return nil
	}
	// 如果目标是数组，将整个cache赋值
	// 如果目标是单个对象，且cache有值，取cache第一个元素赋值，这时候实际上cache最多只有一个元素
	if dataInfo.DataIsSlice {
		dataInfo.Target.Set(cache)
	} else if cache.Len() > 0 {
		dataInfo.Target.Set(cache.Index(0))
	}
	return nil
}

// 用于控制列表处理过程
type FlowControl int

const (
	// go on processing
	Forward FlowControl = 1 + iota
	// continue loop
	Continue
	// break loop
	Break
)
