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
	"testing"
	"time"
)

type TestUnmarshalModel struct {
	FileName           string    `xm:"zip:FileName"`
	FileSize           int       `xm:"zip:FileSize"`
	FileSizeCompressed int       `xm:"zip:FileSizeCompressed"`
	FileDate           time.Time `xm:"zip:FileModified format='2006-01-02 15:03:04' timezone='Asia/Shanghai'"`
	FileComment        string    `xm:"zip:FileComment"`

	Name   string    `xm:"excel://sheet[0]/row[3:]/col[0] pattern='\\w+'"`
	Date   time.Time `xm:"excel://sheet[0]/row[3:]/col[1] pattern='\\d{8}' format='20060102' timezone='Asia/Shanghai'"`
	Open   int       `xm:"excel://sheet[0]/row[3:]/col[4] pattern='\\d+'"`
	High   int       `xm:"excel://sheet[0]/row[3:]/col[5] pattern='\\d+'"`
	Low    int       `xm:"excel://sheet[0]/row[3:]/col[6] pattern='\\d+'"`
	Close  int       `xm:"excel://sheet[0]/row[3:]/col[7] pattern='\\d+'"` //收盘价
	Close2 int       `xm:"excel://sheet[0]/row[3:]/col[8] pattern='\\d+'"` //结算价
	Vol    int       `xm:"excel://sheet[0]/row[3:]/col[11] pattern='\\d+'"`
	Amount int       `xm:"excel://sheet[0]/row[3:]/col[12] pattern='\\d+'"`
	Cang   int       `xm:"excel://sheet[0]/row[3:]/col[13] pattern='\\d+'"`
}

func TestUnmarshal(t *testing.T) {
	reader, e := os.Open("/Users/jiangchangqiang/Desktop/MarketData_Year_2019.zip")
	assert.NoError(t, e)
	unmarshaler := Unmarshaler{
		Charset: "gbk",
		FileFilters: []FileFilter{
			func(fileIndex int, file *zip.File) bool {
				println(">>>>", fileIndex, file.Name)
				return true
			},
		},
		DataLoader: base.DataLoader{
			ItemFilters: []base.ItemFilter{
				func(item interface{}, vars *base.Vars) (flow base.FlowControl, deep int) {
					// check item type
					data, ok := item.(*TestUnmarshalModel)
					if !ok {
						return base.Forward, 0
					}
					// validate item
					//if data.F3 == "118" {
					//	return base.Continue, 1
					//}
					// process item
					println(fmt.Sprintf("%+v", data))
					return base.Forward, 0
				},
			},
		},
	}
	unmarshaler.Unmarshal(reader, &[]*TestUnmarshalModel{})
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
