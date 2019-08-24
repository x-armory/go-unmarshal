package excel

import (
	"errors"
	"github.com/x-armory/go-unmarshal/base"
	"regexp"
	"strconv"
)

const (
	DefaultMaxSheet = 100
	DefaultMaxRow   = 99999
)

// 用于解析base.FieldTag.Path
var posIndexFindReg = regexp.MustCompile(`sheet\[(\w*)(:(\w*))?\]/row\[(\w*)(:(\w*))?\]/col\[(\d+)\]`)

// 用于记录每个field的变量取值范围
type FieldTag struct {
	*base.FieldTag
	SheetStart int
	SheetEnd   int
	RowStart   int
	RowEnd     int
	Col        int
}

// 解析base.FieldTag.Path中的所有变量取值范围，用于遍历提取数据
// tag示例: xm:"xls://sheet[10:20]/row[3:10]/col[3] pattern='\d{4}-\d{2}-\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"
// tag示例: xm:"xls://sheet[]/row[3:]/col[3] pattern='\d{4}-\d{2}-\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"
// tag示例: xm:"xls://sheet[]/row[:4]/col[3]"
// tag示例: xm:"xls://sheet[0]/row[]/col[3]"
// 其中sheet/row下标变量格式为[n:m]或者[n]，n或m为数字是表示常量，非数字表示变量，变量出现在左侧表示(0..x)，出现在右侧表示(x..无穷)，没有冒号表示(0..无穷)
// col下标必须是常量
// 变量默认范围为[0-9999]，可以指定其他值
func GetExcelTags(baseTags *map[int]*base.FieldTag) (map[int]*FieldTag, error) {
	var result = map[int]*FieldTag{}
	allRange := &FieldTag{}
	i := 1
	for _, tag := range *baseTags {
		if tag.Schema != "excel" {
			continue
		}
		index := posIndexFindReg.FindStringSubmatch(tag.Path)
		if len(index) == 0 {
			return nil, errors.New("bad excel tag path")
		}
		excelTag := &FieldTag{FieldTag: tag}
		setRangeIndex(&excelTag.SheetStart, &excelTag.SheetEnd, index[1], index[2] != "", index[3], DefaultMaxSheet)
		setRangeIndex(&excelTag.RowStart, &excelTag.RowEnd, index[4], index[5] != "", index[6], DefaultMaxRow)
		excelTag.Col, _ = strconv.Atoi(index[7])
		result[tag.Id] = excelTag
		if i == 0 {
			allRange.SheetStart = excelTag.SheetStart
			allRange.SheetEnd = excelTag.SheetEnd
			allRange.RowStart = excelTag.RowStart
			allRange.RowEnd = excelTag.RowEnd
		} else {
			allRange.SheetStart = min(allRange.SheetStart, excelTag.SheetStart)
			allRange.SheetEnd = max(allRange.SheetEnd, excelTag.SheetEnd)
			allRange.RowStart = min(allRange.RowStart, excelTag.RowStart)
			allRange.RowEnd = max(allRange.RowEnd, excelTag.RowEnd)
		}
		i++
	}
	result[-1] = allRange

	return result, nil
}

// 根据配置设置改字段的变量取值范围
// start, end 变量取值范围
// startStr 配置的开始取值，为空表示没配置，默认为0
// hasEndStr，是否配置了结束取值
// endStr，配置的结束取值，没配置默认用defaultEnd
// defaultEnd，默认结束取值
func setRangeIndex(start *int, end *int, startStr string, hasEndStr bool, endStr string, defaultEnd int) {
	if startStr == "" && endStr == "" {
		*start = 0
		*end = defaultEnd
	} else {
		*start, _ = strconv.Atoi(startStr)
		if !hasEndStr {
			*end = *start
		} else if v, e := strconv.Atoi(endStr); e == nil {
			*end = v
		} else {
			*end = defaultEnd
		}
		if *end < *start {
			*end = *start
		}
	}
}

func min(a int, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}
func max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}
