package excel

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/x-armory/go-unmarshal/base"
	"regexp"
	"testing"
)

type GetExcelTagsModel struct {
	F1 string `xm:"excel://sheet[5:]/row[3:x]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"`
	F2 string `xm:"excel://sheet[6:9]/row[:]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}'  format='2006-01-02' timezone='Asia/Shanghai'"`
	F3 string `xm:"excel://sheet[:12]/row[:8]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"`
	F4 string `xm:"excel://sheet[5]/row[10:8]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"`
	F5 string `xm:"excel://sheet[]/row[10:8]/col[6] pattern='\\d{4}-\\d{2}-\\d{2}'  format='2006-01-02' timezone='Asia/Shanghai'"`
}

func TestGetExcelTags(t *testing.T) {
	m := GetExcelTagsModel{}
	info, err := base.GetDataInfo(&m)
	assert.NoError(t, err)

	tag, err := GetExcelTags(info.BaseTags)
	assert.NoError(t, err)
	i := -1
	for true {
		if v, ok := tag[i]; ok {
			fmt.Printf("%v\t%+v\n\t%+v\n\n", i, v, v.FieldTag)
		} else {
			break
		}
		i++
	}
}

func TestPosSheetIndexFindReg(t *testing.T) {
	printPosIndex(posIndexFindReg, "sheet[2]/row[3]/col[5]")
	posIndex := printPosIndex(posIndexFindReg, "sheet[]/row[3]/col[5]")
	println(len(posIndex))
	index := printPosIndex(posIndexFindReg, "sheet[2:6]/row[3:3]/col[5]")
	println("sheet", index[1], "-", index[3])
	println("row", index[4], "-", index[6])
	println("col", index[7])
}

func printPosIndex(reg *regexp.Regexp, str string) []string {
	println("[Index]", "reg is", reg.String(), "string is", str)
	submatch := reg.FindStringSubmatch(str)
	for i, m := range submatch {
		fmt.Printf("%d(%s)\n", i, m)
	}
	return submatch
}
