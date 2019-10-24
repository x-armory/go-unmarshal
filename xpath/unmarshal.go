package xpath

import (
	"errors"
	"github.com/x-armory/go-unmarshal/base"
	"gopkg.in/xmlpath.v2"
	"io"
)

// xml文档反序列化；
// 支持的注释格式包括：
// xm:"xpath://*[@id='content_right']/div/div[VarName[n:m]]/a[1]"；
// 其中下标n和m可以为空，缺省最小值为0，缺省最大值为999999
// 按变量定义顺序嵌套循环读取数据；
// 默认遇到空行退出循环
// data数据类型支持：*Obj *[]Obj *[]*Obj；
type Unmarshaler struct {
	base.DataLoader
}

func (m *Unmarshaler) Unmarshal(r io.Reader, data interface{}) error {
	var rt = *m
	// open doc
	doc, e := GetDoc(r)
	if e != nil {
		return e
	}
	// setup
	rt.DataLoader.Data = data
	if rt.ReadValueFunc == nil {
		rt.ReadValueFunc = make(map[string]base.FieldTagReadValueFunc)
	}
	rt.ReadValueFunc["xpath"] = func(fieldTag *base.FieldTag, vars *base.Vars) (v string, err error) {
		return GetValue(doc, fieldTag.PathFilled)
	}
	// load data
	return rt.Load()
}

func GetDoc(r io.Reader) (*xmlpath.Node, error) {
	node, e := xmlpath.ParseHTML(r)
	if e != nil {
		return nil, e
	}
	return node, nil
}

func GetValue(doc *xmlpath.Node, path string) (string, error) {
	xpath, e := xmlpath.Compile(path)
	if e != nil {
		return "", e
	}
	s, ok := xpath.String(doc)
	if ok {
		return s, nil
	} else {
		return "", errors.New("xpath " + path + " not matched")
	}
}
