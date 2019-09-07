package xls

import (
	"fmt"
	"github.com/shakinm/xlsReader/xls"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
	"time"
)

/**
 * 性能好，可以读取所有字符串，只能输入本地文件路径，不接收Reader
 * 不支持日期、公式
 * 因为程序生成的xls内容几本都是字符串内容，暂时不需要支持日期
 *
 * xls的日期有两种格式：
 * https://support.microsoft.com/zh-cn/help/214330/differences-between-the-1900-and-the-1904-date-system-in-excel
 */

func Test_shakinm_xlsReader_1(t *testing.T) {
	workbook, err := xls.OpenFile("/Users/jiangchangqiang/Desktop/MarketData_Year_2019/所内合约行情报表2019.7.xls")
	assert.NoError(t, err)

	sheet, err := workbook.GetSheet(0)
	assert.NoError(t, err)

	row, err := sheet.GetRow(3)
	assert.NoError(t, err)

	c, err := row.GetCol(0)
	assert.NoError(t, err)
	s := c.GetString()

	println(s)
}

func Test_shakinm_xlsReader_2(t *testing.T) {
	workbook, err := xls.OpenFile("/Users/jiangchangqiang/Desktop/test1.xls")
	assert.NoError(t, err)

	sheet, err := workbook.GetSheet(0)
	assert.NoError(t, err)

	row, err := sheet.GetRow(0)
	assert.NoError(t, err)

	{
		c, err := row.GetCol(4)
		assert.NoError(t, err)
		s := c.GetString()
		println(s)
		t := c.GetType()
		println(t)
	}
	{
		c, err := row.GetCol(5)
		assert.NoError(t, err)
		s := c.GetString()
		println(s)
		f, err := strconv.ParseFloat(s, 64)
		assert.NoError(t, err)

		i := int64(f * 24 * 60 * 60)

		date1900, _ := time.Parse("2006-01-02", "1900-01-01")
		date1970, _ := time.Parse("2006-01-02", "1970-01-01")
		unix := time.Unix(i-(date1970.Unix()-date1900.Unix()), 0)
		fmt.Printf("%v", unix)
	}
}

func Test_shakinm_xlsReader_3(t *testing.T) {
	workbook, err := xls.OpenFile("/Users/jiangchangqiang/Desktop/MarketData_Year_2019/所内合约行情报表2019.7.xls")

	if err != nil {
		log.Panic(err.Error())
	}

	// Кол-во листов в книге
	// Number of sheets in the workbook
	//
	// for i := 0; i <= workbook.GetNumberSheets()-1; i++ {}

	fmt.Println(workbook.GetNumberSheets())

	sheet, err := workbook.GetSheet(0)

	if err != nil {
		log.Panic(err.Error())
	}

	// Имя листа
	// Print sheet name
	println(sheet.GetName())

	// Вывести кол-во строк в листе
	// Print the number of rows in the sheet
	println(sheet.GetNumberRows())

	for i := 0; i <= sheet.GetNumberRows(); i++ {
		if row, err := sheet.GetRow(i); err == nil {
			if cell, err := row.GetCol(1); err == nil {

				// Значение ячейки, тип строка
				// Cell value, string type
				fmt.Println(cell.GetString())

				//fmt.Println(cell.GetInt64())
				//fmt.Println(cell.GetFloat64())

				// Тип ячейки (записи)
				// Cell type (records)
				fmt.Println(cell.GetType())

				// Получение отформатированной строки, например для ячеек с датой или проценты
				// Receiving a formatted string, for example, for cells with a date or a percentage
				xfIndex := cell.GetXFIndex()
				formatIndex := workbook.GetXFbyIndex(xfIndex)
				format := workbook.GetFormatByIndex(formatIndex.GetFormatIndex())
				fmt.Println(format.GetFormatString(cell))

			}

		}
	}
}
