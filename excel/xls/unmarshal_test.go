package xls

import (
	"fmt"
	"github.com/extrame/xls"
	"github.com/stretchr/testify/assert"
	"github.com/x-armory/go-unmarshal/base"
	"testing"
)

type TestXlsUnmarshalModel struct {
	FileName string `xm:"zip://filename"`
	F1       string `xm:"excel://sheet[]/row[]/col[0]"`
	F2       string `xm:"excel://sheet[]/row[]/col[1]"`
	F3       string `xm:"excel://sheet[]/row[]/col[2]"`
	F4       string `xm:"excel://sheet[]/row[]/col[3]"`
}

func TestXls(t *testing.T) {
	book, e := xls.Open("/Users/jiangchangqiang/Desktop/MarketData_Year_2019/所内合约行情报表2019.7.xls", "gbk")
	assert.NoError(t, e)

	for i := 0; i < 10; i++ {
		row := book.GetSheet(0).Row(i)
		print(i)
		for c := row.FirstCol(); c <= row.LastCol(); c++ {
			print("\t", row.Col(c))
		}
		println()
	}
}

func TestXlsUnmarshal(t *testing.T) {
	book, e := xls.Open("/Users/jiangchangqiang/Desktop/test1.xls", "utf-8")
	assert.NoError(t, e)
	var data []TestXlsUnmarshalModel
	e = Unmarshal(book, &data, false, func(item interface{}) base.FlowControl {
		fmt.Printf("%v\n", item)
		return base.Forward
	})
	assert.NoError(t, e)

	fmt.Printf("%v", data)
}
