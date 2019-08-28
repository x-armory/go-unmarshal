package zip

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/x-armory/go-unmarshal/base"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"
)

type TestUnmarshalModel struct {
	Name   string    `xm:"excel://sheet[0]/row[3:]/col[0]"`
	Date   time.Time `xm:"excel://sheet[0]/row[3:]/col[1] format='20060102' timezone='Asia/Shanghai'"`
	Open   float64   `xm:"excel://sheet[0]/row[3:]/col[4]"`
	High   float64   `xm:"excel://sheet[0]/row[3:]/col[5]"`
	Low    float64   `xm:"excel://sheet[0]/row[3:]/col[6]"`
	Close  float64   `xm:"excel://sheet[0]/row[3:]/col[7]"` //收盘价
	Close2 float64   `xm:"excel://sheet[0]/row[3:]/col[8]"` //结算价
	Vol    float64   `xm:"excel://sheet[0]/row[3:]/col[11]"`
	Amount float64   `xm:"excel://sheet[0]/row[3:]/col[12]"`
	Cang   float64   `xm:"excel://sheet[0]/row[3:]/col[13]"`
}

func TestUnmarshal(t *testing.T) {
	reader, e := GetZipReaderFromLocal("/Users/jiangchangqiang/Desktop/MarketData_Year_2019.zip")
	assert.NoError(t, e)

	var lastSyncDate = time.Now().AddDate(0, -1, 0)
	var data []TestUnmarshalModel
	var name string
	var fileNameMonthReg = regexp.MustCompile("(\\d+)[.](\\d+)[.]xls$")
	var row int

	assert.NoError(t,
		Unmarshal(reader, "gbk", func(fileName string) base.FlowControl {
			// 检查文件名中的日期部分，如果找到了年月，说明是今年数据，从上次同步过的日期开始
			match := fileNameMonthReg.FindStringSubmatch(fileName)
			if len(match) == 3 {
				year, _ := strconv.Atoi(match[1])
				month, _ := strconv.Atoi(match[2])
				if year*100+month < lastSyncDate.Year()*100+int(lastSyncDate.Month()) {
					println("ignore", fileName)
					return base.Continue
				}
				println("process", fileName)
				row = 0
			}
			return base.Forward
		}, &data, false, func(item interface{}) base.FlowControl {
			// 数据标题后面可能有些空行，跳过空行
			row++
			if item == nil {
				if row < 5 {
					return base.Continue
				} else {
					return base.Break
				}
			}
			return base.Forward
		}, func(item interface{}) base.FlowControl {
			// 补全数据
			model := item.(*TestUnmarshalModel)
			if model.Name != "" {
				name = model.Name
			} else {
				model.Name = name
			}
			// 从7月开始读取数据
			if model.Date.Month() > 6 {
				return base.Forward
			} else {
				return base.Continue
			}
		}, func(item interface{}) base.FlowControl {
			// 处理韩数据
			model := item.(*TestUnmarshalModel)
			fmt.Printf("save %d\t%+v\n", row, model)
			return base.Forward
		}))
	println(len(data))
}

func TestGetZipReaderFromLocal(t *testing.T) {
	reader, e := GetZipReaderFromLocal("/Users/jiangchangqiang/Desktop/MarketData_Year_2018.zip")
	assert.NoError(t, e)
	assert.Equal(t, true, reader != nil)
}

func TestGetFileExtName(t *testing.T) {
	file := "file.ext1.ext"
	reg := regexp.MustCompile("[.](\\w+)$")
	submatch := reg.FindStringSubmatch(file)
	fmt.Printf("%v", submatch[1])
}

func TestCharset(t *testing.T) {
	file, e := os.Open("/Users/jiangchangqiang/Desktop/MarketData_Year_2019.zip")
	assert.NoError(t, e)
	i := determineEncoding(file)
	fmt.Printf("%v\n", i)

	//reader := transform.NewReader(file, i.NewDecoder())
	info, e := file.Stat()
	assert.NoError(t, e)

	newReader, e := zip.NewReader(file, info.Size())
	assert.NoError(t, e)

	for _, f := range newReader.File {
		closer, e := f.Open()
		assert.NoError(t, e)
		fileEncoding := determineEncoding(closer)

		s, e := fileEncoding.NewDecoder().String(f.Name)
		assert.NoError(t, e)

		e2, name := charset.Lookup("gbk")
		fmt.Printf("%v\t%v\n", e2, name)

		i2, e := simplifiedchinese.GBK.NewDecoder().String(f.Name)
		assert.NoError(t, e)

		println(f.Name)
		println(s)
		println(i2)
		fmt.Printf("%v\n", fileEncoding)
	}
}
func determineEncoding(r io.Reader) encoding.Encoding {
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
