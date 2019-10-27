package xls

import (
	"fmt"
	"github.com/extrame/xls"
	"github.com/stretchr/testify/assert"
	"github.com/x-armory/go-unmarshal/base"
	"os"
	"testing"
)

type TestXlsUnmarshalModel struct {
	F1 string `xm:"excel:sheet[:]/row[3:]/col[0]"`
	F2 string `xm:"excel:sheet[:]/row[3:]/col[1]"`
	F3 string `xm:"excel:sheet[:]/row[3:]/col[2]"`
	F4 string `xm:"excel:sheet[:]/row[3:]/col[3]"`
}

func TestXls2(t *testing.T) {
	book, e := xls.Open("/Users/jiangchangqiang/Desktop/MarketData_Year_2019/所内合约行情报表2019.7.xls", "gbk")
	assert.NoError(t, e)
	res := book.ReadAllCells(50000)
	strings := res[20630]
	println(strings)
}
func TestXls(t *testing.T) {
	book, e := xls.Open("/Users/jiangchangqiang/Desktop/MarketData_Year_2019/所内合约行情报表2019.7.xls", "gbk")
	assert.NoError(t, e)

	r1 := book.GetSheet(0).
		Row(20630)
	r2 := book.GetSheet(0).
		Row(4)

	fmt.Printf("%v\n%v",
		r1.Col(0),
		r2.Col(0))
}

func TestXlsUnmarshal(t *testing.T) {
	file, e := os.Open("/Users/jiangchangqiang/Desktop/test1.xls")
	assert.NoError(t, e)
	xlsUnmarshaler := &Unmarshaler{
		DataLoader: base.DataLoader{
			ItemFilters: []base.ItemFilter{
				func(item interface{}, vars *base.Vars) (flow base.FlowControl, deep int) {
					// check item type
					data, ok := item.(*TestXlsUnmarshalModel)
					if !ok {
						return base.Forward, 0
					}
					// validate item
					if data.F3 == "118" {
						return base.Continue, 1
					}
					// process item
					fmt.Printf("%+v\n", data)
					return base.Forward, 0
				},
			},
		},
	}

	var data []TestXlsUnmarshalModel
	e = xlsUnmarshaler.Unmarshal(file, &data)
	assert.NoError(t, e)
	assert.Equal(t, 0, len(data))
}
