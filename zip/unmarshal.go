package zip

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/x-armory/go-unmarshal/base"
	"github.com/x-armory/go-unmarshal/csv"
	"github.com/x-armory/go-unmarshal/excel/xls"
	"github.com/x-armory/go-unmarshal/xpath"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// zip反序列化，按文件名顺序读取文件，根据文件后缀名选择反序列化工具；
// 目前支持excel(xls)、xml(html/htm/xml/xhtml)两种格式；
// Charset，默认utf-8；
// varOrder，变量嵌套顺序，默认随机，excel文档固定按 sheet->row->col遍历，指定无效；
// writeDate，写数据开关，因为zip包数据流通常非常大，暂时设置无效，统一不写数据；
// FileFilters，文件过滤器，返回false时跳过文件，可用于截取文件中的变量，并在后续itemFilters中设置到item里，比如文件名中可能存在的分类、日期等；
// itemFilters，目标对象元素过滤器，用于校验、处理元素，并可控制文件反序列化流程；
// 目标对象类型可使用tag xm:"zip:VarName pattern='' format='' timezone=''"获取文件相关信息，支持的var变量包括：
// FileName：string，文件名；
// Comment：string，文件备注；
// Modified：time.Time 文件修改日期，固定格式：2006-01-02 15:03:04；
// CompressedSize64：int64 压缩后大小；
// UncompressedSize64：int64 解压后大小；
type Unmarshaler struct {
	Charset        string
	charsetDecoder *encoding.Decoder
	FileFilters    []FileFilter
	base.DataLoader
}
type FileFilter func(fileIndex int, file *zip.File) bool

func (m *Unmarshaler) Unmarshal(r io.Reader, data interface{}) error {
	var rt = *m
	if rt.Charset == "" {
		rt.Charset = "utf-8"
	}
	switch strings.ToLower(rt.Charset) {
	case "gbk":
		rt.charsetDecoder = simplifiedchinese.GBK.NewDecoder()
	}
	// zip暂时不允许写数据，因为数据量通常很大
	rt.DataLoader.WriteData = false
	//读取zip
	doc, e := GetDoc(r, rt.Charset)
	if e != nil {
		return e
	}
	// 准备反序列化工具，在定制itemFilters之前设置zip注解处理filter
	rt.Data = data
	if rt.ReadValueFunc == nil {
		rt.ReadValueFunc = make(map[string]base.FieldTagReadValueFunc)
	}
	// 按顺序反序列化所有文件
rootLoop:
	for i, file := range doc.File {
		// 过滤、校验文件，可用于筛选文件类型，获取文件名变量等
		for _, filter := range rt.FileFilters {
			if !filter(i, file) {
				continue rootLoop
			}
		}
		// 根据扩展名准备文件反序列化工具
		loader := rt.DataLoader
		loader.ReadValueFunc["zip"] = func(fieldTag *base.FieldTag, vars *base.Vars) (v string, err error) {
			return GetValue(file, fieldTag.PathFilled)
		}
		var unmarshal base.Unmarshaler
		ext := filepath.Ext(file.Name)
		switch ext {
		default:
			continue
		case ".csv", ".txt":
			unmarshal = &csv.Unmarshaler{DataLoader: loader}
		case ".xls":
			unmarshal = &xls.Unmarshaler{Charset: rt.Charset, DataLoader: loader}
		case ".html", ".htm", ".xhtml", ".xml":
			unmarshal = &xpath.Unmarshaler{DataLoader: loader}
		}
		// 执行反序列化
		if e := func() error {
			var reader io.Reader
			fileReader, e := file.Open()
			if e != nil {
				return e
			}
			defer fileReader.Close()
			if rt.charsetDecoder == nil {
				reader = fileReader
			} else {
				reader = transform.NewReader(fileReader, rt.charsetDecoder)
			}
			if e := unmarshal.Unmarshal(reader, data); e != nil {
				return e
			}
			return nil
		}(); e != nil {
			println("[WARN]", e.Error())
		}
	}
	return nil
}

// 打开zip，解码文件名，排序文件
func GetDoc(r io.Reader, charsetName string) (*zip.Reader, error) {
	if charsetName == "" {
		charsetName = "utf-8"
	}
	// load Charset
	encoding, _ := charset.Lookup(charsetName)
	if encoding == nil {
		return nil, errors.New("Charset " + charsetName + " not supported")
	}
	// read content
	bts, e := ioutil.ReadAll(r)
	if e != nil {
		return nil, e
	}
	// open zip
	reader, e := zip.NewReader(bytes.NewReader(bts), int64(len(bts)))
	if e != nil {
		return nil, e
	}
	// validate zip
	if len(reader.File) == 0 {
		return nil, errors.New("no file found in zip")
	}
	// decode file name
	if strings.ToLower(charsetName) != "utf-8" {
		for _, file := range reader.File {
			file.Name, e = encoding.NewDecoder().String(file.Name)
			if e != nil {
				return nil, e
			}
		}
	}
	// sort file by name
	sort.Slice(reader.File, func(i, j int) bool {
		return reader.File[i].Name <= reader.File[j].Name
	})

	return reader, nil
}

func GetValue(doc *zip.File, path string) (string, error) {
	switch path {
	case "FileName":
		return doc.Name, nil
	case "FileComment":
		return doc.Comment, nil
	case "FileModified":
		return doc.Modified.Format("2006-01-02 15:03:04"), nil
	case "FileSizeCompressed":
		return strconv.FormatUint(doc.CompressedSize64, 10), nil
	case "FileSize":
		return strconv.FormatUint(doc.UncompressedSize64, 10), nil
	}
	return "", nil
}
