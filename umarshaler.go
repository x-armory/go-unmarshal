package go_unmarshal

import (
	"github.com/extrame/xls"
	"github.com/x-armory/go-unmarshal/base"
	"gopkg.in/xmlpath.v2"
)

// 根据input类型选择反序列化工具
// 支持zip、xls/xlsx、csv、xmlpath
func Unmarshal(input interface{}, data interface{}, writeData bool, itemFilter ...base.ItemFilter) error {
	switch input.(type) {
	case *xls.WorkBook:
		return nil
	case *xmlpath.Node:
		return nil
	}
}
