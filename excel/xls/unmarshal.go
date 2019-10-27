package xls

import (
	"errors"
	"fmt"
	"github.com/shakinm/xlsReader/xls"
	"github.com/x-armory/go-unmarshal/base"
	"github.com/x-armory/go-unmarshal/excel"
	"io"
	"math/rand"
	"os"
	"path"
)

// xls格式excel反序列化；
// 支持的注释格式包括：
// xm:"excel:sheet[1:2]/row[1:30]/col[1]"；
// xm:"excel:sheet[1:2]/row[1:]/col[1]"；
// xm:"excel:sheet[1:2]/row[:30]/col[1]"；
// xm:"excel:sheet[1:2]/row[:]/col[1]"；
// 其中sheet、row和col下标从0开始；
// 下标缺省最小值为0，缺省最大值为999999
// 按sheet->row->col顺序嵌套循环读取数据；
// 默认遇到空行退出循环
// data数据类型支持：*Obj *[]Obj *[]*Obj；
type Unmarshaler struct {
	base.DataLoader
}

func (m *Unmarshaler) Unmarshal(r io.Reader, data interface{}) error {
	var rt = *m
	// open excel doc
	doc, e := GetDoc(r)
	if e != nil {
		return e
	}
	// setup
	rt.DataLoader.Data = data
	// 固定变量嵌套循环顺序
	rt.VarOrder = []string{"sheet", "row", "col"}
	if rt.ExitNoDataTimes <= 0 || rt.ExitNoDataTimes > 10 {
		rt.ExitNoDataTimes = 3
	}
	if rt.ReadValueFunc == nil {
		rt.ReadValueFunc = make(map[string]base.FieldTagReadValueFunc)
	}
	rt.ReadValueFunc["excel"] = func(fieldTag *base.FieldTag, vars *base.Vars) (v string, err error) {
		return GetValue(doc, fieldTag.PathFilled)
	}
	// load data
	return rt.Load()
}

func GetDoc(r io.Reader) (*xls.Workbook, error) {
	dir := path.Join(os.TempDir(), fmt.Sprintf("xarmory/go-unmarshal/xls/%d", rand.Int63()))
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return nil, err
	}
	file := path.Join(dir, "tmp")
	tmp, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(tmp, r)
	if err != nil {
		return nil, err
	}

	workbook, err := xls.OpenFile(file)
	if err != nil {
		return nil, err
	}
	return &workbook, nil
}

func GetValue(doc *xls.Workbook, path string) (s string, err error) {
	sheet, row, col, err := excel.GetVar(path)
	if err != nil {
		return "", err
	}

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("get sheet failed, %v, sheet=%d, row=%d, col=%d", e, sheet, row, col))
		}
	}()

	sheetDoc, err := doc.GetSheet(sheet)
	if err != nil {
		return "", errors.New(fmt.Sprintf("get sheet failed, %s, sheet=%d, row=%d, col=%d", err.Error(), sheet, row, col))
	}
	rowDoc, err := sheetDoc.GetRow(row)
	if err != nil {
		return "", errors.New(fmt.Sprintf("get row failed, %s, sheet=%d, row=%d, col=%d", err.Error(), sheet, row, col))
	}
	c, err := rowDoc.GetCol(col)
	if err != nil {
		return "", errors.New(fmt.Sprintf("get col failed, %s, sheet=%d, row=%d, col=%d", err.Error(), sheet, row, col))
	}
	return c.GetString(), nil
}
