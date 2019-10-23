package base

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"regexp"
	"testing"
)

type TestGetFieldTagsModel struct {
	F1 string `xm:"excel:sheet[5:7]/row[3:5]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"`
	F2 string `xm:"excel:sheet[6]/row[:4:2]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' patternIdx='3' format='2006-01-02' timezone='Asia/Shanghai'"`
	F3 string `xm:"excel:sheet[6]/row[:4:2]/col[3] pattern='\\d+' patternIdx='0' format='2006-01-02' timezone='Asia/Shanghai'"`
}

func TestGetFieldTags(t *testing.T) {
	model := TestGetFieldTagsModel{}
	tp := reflect.TypeOf(model)
	tags, e := GetFieldTags(tp, "xm", nil)
	fmt.Printf("%v\n%v\n\n", e, tags)
	for k, v := range *tags {
		bts, _ := json.Marshal(v)
		fmt.Printf("%v\t%+v\n", k, string(bts))
	}
}

func TestFieldTag_Fill(t *testing.T) {
	model := TestGetFieldTagsModel{}
	tp := reflect.TypeOf(model)
	tags, e := GetFieldTags(tp, "xm", nil)
	assert.NoError(t, e)
	tags = tags.Filter(func(tag *FieldTag) bool {
		return tag.Schema == "excel"
	})
	vars := tags.MergeVars().List("sheet", "row", "col")
	for vars.Reset(); vars.IsValid(); vars.Next() {
		tags.SetValues(vars)
		bts, _ := json.Marshal(tags)
		fmt.Printf("%+v\n", string(bts))
	}
}

func TestFieldTag_Parse(t *testing.T) {
	model := TestGetFieldTagsModel{}
	tp := reflect.TypeOf(model)
	tags, _ := GetFieldTags(tp, "xm", nil)
	tag := (*tags)[2]
	value, e := tag.Parse("1,2,3,4")
	fmt.Printf("%v %v\n", value, e)
}

func TestReg(t *testing.T) {
	allString := regexp.MustCompile("\\d+").FindAllString("3:4", -1)
	println(len(allString))
}
