package base

import (
	"fmt"
	"reflect"
	"testing"
)

type TestGetFieldTagsModel struct {
	F1 string `xm:"xls://sheet[5:7]/row[3:x]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"`
	F2 string `xm:"xls://sheet[6]/row[:]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' patternIdx='3' format='2006-01-02' timezone='Asia/Shanghai'"`
}

func TestGetFieldTags(t *testing.T) {
	model := TestGetFieldTagsModel{}
	tp := reflect.TypeOf(model)
	tags, e := GetFieldTags(tp, "xm")
	fmt.Printf("%v\n%v\n\n", e, tags)
	for k, v := range *tags {
		fmt.Printf("%v\t%+v\n", k, v)
	}
}
