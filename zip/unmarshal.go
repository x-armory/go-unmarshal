package zip

import (
	"archive/zip"
	"errors"
	"github.com/x-armory/go-unmarshal/base"
	"github.com/x-armory/go-unmarshal/excel/xls"
	"golang.org/x/net/html/charset"
	"os"
	"regexp"
	"sort"
)

func GetZipReaderFromLocal(url string) (*zip.Reader, error) {
	file, e := os.Open(url)
	if e != nil {
		return nil, e
	}
	return func() (*zip.Reader, error) {
		info, e := file.Stat()
		if e != nil {
			return nil, e
		}

		return zip.NewReader(file, info.Size())
	}()
}

//func GetZipReaderFromReader(r *io.Reader) (*zip.Reader, error) {
//
//}

// 返回true继续处理文件
// 返回false跳过
type ItemFilter func(fileName string) base.FlowControl

var fileExtNameReg = regexp.MustCompile("[.](\\w+)$")

func Unmarshal(doc *zip.Reader, charsetName string, fileFilter ItemFilter, data interface{}, writeDate bool, filters ...base.ItemFilter) error {
	encoding, _ := charset.Lookup("gbk")
	if encoding == nil {
		return errors.New("charset " + charsetName + " not supported")
	}
	sort.Slice(doc.File, func(i, j int) bool {
		return doc.File[i].Name <= doc.File[j].Name
	})
	for _, itemFile := range doc.File {
		var fileExtName = ""
		var err error
		fileName := itemFile.Name
		fileName, err = encoding.NewDecoder().String(fileName)
		if err != nil {
			return err
		}
		match := fileExtNameReg.FindStringSubmatch(fileName)
		if len(match) != 2 {
			continue
		}
		fileExtName = match[1]
		fileFlow := fileFilter(fileName)
		if fileFlow == base.Continue {
			continue
		} else if fileFlow == base.Break {
			break
		}
		switch fileExtName {
		default:
			continue
		case "xls":
			itemReader, err := itemFile.Open()
			if err != nil {
				return err
			}
			book, err := xls.GetWorkBook(itemReader, charsetName)
			if err != nil {
				return err
			}
			err = xls.Unmarshal(book, data, writeDate, filters...)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
