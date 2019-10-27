package csv

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/x-armory/go-unmarshal/base"
	"golang.org/x/text/encoding"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// ColCount不为0时，用于过滤长度不匹配的行
type Unmarshaler struct {
	base.DataLoader
	Encoding     encoding.Encoding
	GetSepFunc   func() string
	RowParseFunc func(string) []string
}

func DefaultGetSepFunc() string {
	return "\t"
}

func (m *Unmarshaler) Unmarshal(r io.Reader, data interface{}) error {
	var rt = *m
	if m.GetSepFunc == nil {
		m.GetSepFunc = DefaultGetSepFunc
	}
	err, newReader := base.TransformReaderEncoding(r, rt.Encoding)
	if err != nil {
		return err
	}
	// open doc
	doc, e := m.GetDoc(newReader)
	if e != nil {
		return e
	}
	// setup
	rt.DataLoader.Data = data
	if rt.ReadValueFunc == nil {
		rt.ReadValueFunc = make(map[string]base.FieldTagReadValueFunc)
	}
	rt.ReadValueFunc["csv"] = func(fieldTag *base.FieldTag, vars *base.Vars) (v string, err error) {
		return GetValue(doc, fieldTag.PathFilled)
	}
	// load data
	return rt.Load()
}

func (m *Unmarshaler) GetDoc(r io.Reader) (re [][]string, err error) {
	reader := bufio.NewReader(r)
	for true {
		line, _, e := reader.ReadLine()
		if e != nil {
			if e != io.EOF {
				err = e
			}
			break
		}
		content := string(line)
		var split []string
		if m.RowParseFunc != nil {
			split = m.RowParseFunc(content)
		} else {
			split = strings.Split(content, m.GetSepFunc())
		}
		re = append(re, split)
	}
	return re, nil
}

func GetValue(doc [][]string, path string) (string, error) {
	row, col, err := GetVar(path)
	if err != nil {
		return "", err
	}
	if row < 0 || row > len(doc)-1 {
		return "", errors.New(fmt.Sprintf("row %d out of range[0-%d]", row, len(doc)-1))
	}
	line := doc[row]
	if col < 0 || col > len(line)-1 {
		return "", errors.New(fmt.Sprintf("col %d out of range[0-%d]", col, len(line)-1))
	}
	return line[col], nil
}

var FindPathVarReg = regexp.MustCompile(`//(row\[)?(\d+)(\])?/(col\[)?(\d+)(\])?`)

func GetVar(pathFilled string) (row int, col int, err error) {
	matches := FindPathVarReg.FindStringSubmatch(pathFilled)
	if len(matches) != 7 {
		return -1, -1, errors.New(pathFilled + " not match " + FindPathVarReg.String())
	}
	row, _ = strconv.Atoi(matches[2])
	col, _ = strconv.Atoi(matches[5])
	return row, col, nil
}
