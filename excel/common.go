package excel

import (
	"errors"
	"github.com/x-armory/go-unmarshal/base"
	"regexp"
	"strconv"
)

// 用于校验base.FieldTag.Path
var FindPathReg = regexp.MustCompile(`sheet\[(\d*)(:(\d*))?\]/row\[(\d*)(:(\d*))?\]/col\[(\d+)\]`)
var FindPathVarReg = regexp.MustCompile(`(sheet\[)?(\d+)(\])?/(row\[)?(\d+)(\])?/(col\[)?(\d+)(\])?`)

// 获取schema是excel、匹配excelPathReg的所有tag
func GetExcelTags(baseTags *base.FieldTagMap) *base.FieldTagMap {
	return baseTags.Filter(func(tag *base.FieldTag) bool {
		return tag.Schema == "excel" && FindPathReg.MatchString(tag.Path)
	})
}

func GetVar(pathFilled string) (sheet int, row int, col int, err error) {
	matches := FindPathVarReg.FindStringSubmatch(pathFilled)
	if len(matches) != 10 {
		return -1, -1, -1, errors.New(pathFilled + " not match " + FindPathVarReg.String())
	}
	sheet, _ = strconv.Atoi(matches[2])
	row, _ = strconv.Atoi(matches[5])
	col, _ = strconv.Atoi(matches[8])
	return sheet, row, col, nil
}
