package base

import (
	"fmt"
	"testing"
	"time"
)

type GetDataLoaderModel struct {
	F1 string     `excel:"sheet[5:7]/row[3:x]/col[3]"`
	F2 time.Time  `excel:"sheet[6]/row[:]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"`
	F3 *time.Time `excel:"sheet[6]/row[:]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"`
	F4 int        `excel:"sheet[6]/row[:]/col[3] pattern='\\d+'"`
}

func TestGetDataLoader(t *testing.T) {
	var dataChan = make(chan ItemAllFieldsValue, 2)
	go func() {
		for i := 0; i < 9; i++ {
			print("set data", i, "\n")
			dataChan <- &map[int]string{
				0: fmt.Sprintf("row[%d]/F1", i),
				1: fmt.Sprintf("201%d-01-01", i),
				2: fmt.Sprintf("2019-0%d-01", i+1),
				3: fmt.Sprintf("row[%d]/F4", i),
			}
		}
		// stop sig
		close(dataChan)
	}()

	var model []GetDataLoaderModel

	loader, e := NewDataLoader("excel", &model, true,
		func(item interface{}) bool {
			tagsModel := item.(*GetDataLoaderModel)
			return tagsModel.F4 <= 3
		},
		func(item interface{}) bool {
			if i, ok := item.(*GetDataLoaderModel); ok {
				fmt.Printf("Got an item : %+v\n", i)
			}
			return true
		})
	if e != nil {
		t.Error(e)
	}

	e = loader.LoadData(dataChan)
	if e != nil {
		t.Error(e)
	}

	fmt.Printf("data: %+v", model)
}
