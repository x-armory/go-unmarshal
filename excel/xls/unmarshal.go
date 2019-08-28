package xls

import (
	"bytes"
	"github.com/extrame/xls"
	"github.com/x-armory/go-unmarshal/base"
	"github.com/x-armory/go-unmarshal/excel"
	"io"
	"io/ioutil"
)

func GetWorkBook(r io.Reader, charset string) (*xls.WorkBook, error) {
	bs, e := ioutil.ReadAll(r)
	if e != nil {
		return nil, e
	}
	return xls.OpenReader(bytes.NewReader(bs), charset)
}

// data support *Obj *[]Obj *[]*Obj
func Unmarshal(doc *xls.WorkBook, data interface{}, writeDate bool, filters ...base.ItemFilter) error {
	loader, err := base.NewDataLoader(nil, data, writeDate, filters...)
	if err != nil {
		return err
	}
	var dataChan = make(base.ItemValueChan)
	// fetch data
	go func() {
		tags, e := excel.GetExcelTags(loader.BaseTags)
		if e != nil {
			err = e
			return
		}
		allRange := tags[-1]
		delete(tags, -1)
		for sheet := allRange.SheetStart; sheet <= allRange.SheetEnd && sheet < doc.NumSheets(); sheet++ {
			var sheetDoc = doc.GetSheet(sheet)
			var sheetMaxRow = sheetDoc.MaxRow
			var sheetMaxRowInt = int(sheetMaxRow)
			for row := allRange.RowStart; true; row++ {
				if !(row <= allRange.RowEnd && row <= sheetMaxRowInt) {
					break
				}
				var rowDoc = sheetDoc.Row(row)
				var itemValues = map[int]string{}
				for field, fieldTag := range tags {
					if field < rowDoc.FirstCol() || field > rowDoc.LastCol() {
						itemValues[fieldTag.Id] = ""
					} else {
						itemValues[fieldTag.Id] = rowDoc.Col(fieldTag.Col)
					}
				}
				dataChan <- &itemValues
				if !loader.DataIsSlice {
					return
				}
			}
		}
		println("END")
	}()
	e := loader.LoadData(dataChan)
	if e != nil {
		return e
	}
	if err != nil {
		return err
	}
	return nil
}
