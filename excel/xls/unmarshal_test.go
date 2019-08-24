package xls

import (
	"fmt"
	"github.com/extrame/xls"
	"github.com/stretchr/testify/assert"
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
	book, e := xls.Open("/Users/jiangchangqiang/Desktop/test1.xls", "utf-8")
	assert.NoError(t, e)
	fmt.Printf("%v", book.GetSheet(0).Row(15).Col(46))
	println(book.NumSheets())
	println(book.GetSheet(0).MaxRow)
	println(book.GetSheet(0).Row(0).LastCol())
	println(book.GetSheet(0).Row(0).FirstCol())
}

func TestXlsUnmarshal(t *testing.T) {
	book, e := xls.Open("/Users/jiangchangqiang/Desktop/test1.xls", "utf-8")
	assert.NoError(t, e)
	var data []TestXlsUnmarshalModel
	e = Unmarshal(book, &data, false, func(item interface{}) bool {
		fmt.Printf("%v\n", item)
		return true
	})
	assert.NoError(t, e)

	fmt.Printf("%v", data)
}
