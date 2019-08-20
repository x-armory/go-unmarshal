package base

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var posSchemaFindReg = regexp.MustCompile("^(\\w+)[:]//(.+)")
var posPatternFindReg = regexp.MustCompile(" +pattern='([^']+)'")
var posPatternIndexFindReg = regexp.MustCompile(" +patternIdx='([^']+)'")
var posFormatFindReg = regexp.MustCompile(" +format='([^']+)'")
var posTimezoneFindReg = regexp.MustCompile(" +timezone='([^']+)'")

type FieldTag struct {
	Id         int            //field index
	Name       string         //field name
	FieldType  reflect.Type   //
	Schema     string         // zip | xls | xlsx | xpath | csv
	Path       string         //position expression; e.g. sheet[x:x]/row[3:x]/col[3]
	Pattern    *regexp.Regexp //find value
	PatternIdx int            //find value
	Format     string         //format value, only for time
	Timezone   *time.Location //default +8, only for time
}

func GetFieldTags(T reflect.Type, tag string) (*map[int]*FieldTag, error) {
	var result = map[int]*FieldTag{}
	var t = T
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fieldNum := t.NumField()
	var lastProtocol = ""
	for i := 0; i < fieldNum; i++ {
		f := t.Field(i)
		posConfig, ok := f.Tag.Lookup(tag)
		if !ok {
			continue
		}
		pos := strings.TrimSpace(posConfig)
		finds := posSchemaFindReg.FindStringSubmatch(pos)
		if len(finds) != 3 {
			return nil, errors.New("bad xm tag " + pos)
		}
		schema := finds[1]
		if lastProtocol == "" {
			lastProtocol = schema
		} else if lastProtocol != schema {
			return nil, errors.New(fmt.Sprintf("multiple schemas not supported, %s, %s ", lastProtocol, schema))
		}

		pos = finds[2]
		pos, patternStr := splitPosConfig(pos, posPatternFindReg, "")
		pos, patternIdxStr := splitPosConfig(pos, posPatternIndexFindReg, "")
		pos, formatStr := splitPosConfig(pos, posFormatFindReg, "")
		pos, timezoneStr := splitPosConfig(pos, posTimezoneFindReg, "Asia/Shanghai")
		location, err := time.LoadLocation(timezoneStr)
		if err != nil {
			return nil, err
		}
		var pattern *regexp.Regexp = nil
		if patternStr != "" {
			pattern = regexp.MustCompile(patternStr)
		}
		patternIdx, _ := strconv.Atoi(patternIdxStr)

		result[i] =
			&FieldTag{
				i,
				f.Name,
				f.Type,
				schema,
				pos,
				pattern,
				patternIdx,
				formatStr,
				location,
			}
	}
	if len(result) == 0 {
		return nil, errors.New("xm tag not found")
	}
	return &result, nil
}

func (tag *FieldTag) Parse(str string) (reflect.Value, error) {
	if tag.Pattern != nil {
		subMatches := tag.Pattern.FindStringSubmatch(str)
		if len(subMatches) > 0 && subMatches[0] != "" {
			str = subMatches[len(subMatches)-1]
		} else {
			str = ""
		}
	}
	if str == "" {
		return reflect.Value{}, nil
	}
	switch tag.FieldType.Kind() {
	case reflect.String:
		return reflect.ValueOf(str), nil
	case reflect.Int:
		i, err := strconv.Atoi(str)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), err
	case reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(f), err
	case reflect.Struct, reflect.Ptr:
		var isTime = tag.FieldType.AssignableTo(reflect.TypeOf(time.Time{}))
		var isTimePtr = tag.FieldType.AssignableTo(reflect.TypeOf(&time.Time{}))
		if !isTime && !isTimePtr {
			return reflect.Value{}, errors.New(fmt.Sprintf("%v not supported", tag.FieldType))
		}
		t, err := time.ParseInLocation(tag.Format, str, tag.Timezone)
		if err != nil {
			return reflect.Value{}, err
		}
		if isTime {
			return reflect.ValueOf(t), err
		} else if isTimePtr {
			return reflect.ValueOf(&t), err
		}

	}
	return reflect.Value{}, nil
}

func splitPosConfig(str string, pattern *regexp.Regexp, defaultV string) (string, string) {
	subMatches := pattern.FindStringSubmatch(str)
	if len(subMatches) == 2 {
		return strings.Replace(str, subMatches[0], "", 1), subMatches[1]
	} else {
		return str, defaultV
	}
}
