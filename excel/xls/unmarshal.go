package xls

import (
	"github.com/extrame/xls"
	"github.com/x-armory/go-unmarshal/base"
	"github.com/x-armory/go-unmarshal/excel"
)

// data support *Obj *[]Obj *[]*Obj
func Unmarshal(doc *xls.WorkBook, data interface{}, writeDate bool, filters ...base.ItemFilter) error {
	loader, err := base.NewDataLoader(nil, data, writeDate, filters...)
	if err != nil {
		return err
	}
	var dataChan = make(base.ItemValueChan, 5)
	// fetch data
	go func() {
		defer close(dataChan)
		tags, e := excel.GetExcelTags(loader.BaseTags)
		if e != nil {
			err = e
			return
		}
		allRange := tags[-1]
		delete(tags, -1)
		for sheet := allRange.SheetStart; sheet <= allRange.SheetEnd && sheet < doc.NumSheets(); sheet++ {
			var sheetDoc = doc.GetSheet(sheet)
			for row := allRange.RowStart; row <= allRange.RowEnd && row <= int(sheetDoc.MaxRow); row++ {
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
