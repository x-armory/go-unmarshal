package xls

import (
	"fmt"
	"github.com/extrame/xls"
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
 * 性能较差，很多自定义格式无法解析
 */

func Test_extrame_xls_1(t *testing.T) {
	book, e := xls.Open("/Users/jiangchangqiang/Desktop/MarketData_Year_2019/所内合约行情报表2019.7.xls", "gbk")
	assert.NoError(t, e)
	res := book.ReadAllCells(50000)
	strings := res[20630]
	println(strings)
}
func Test_extrame_xls_2(t *testing.T) {
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
