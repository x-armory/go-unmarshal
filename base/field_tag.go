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

// path表达式，格式为：schema://...
var posSchemaFindReg = regexp.MustCompile("^(\\w+)[:](.+)")

// 取值表达式，格式为：pattern='regexp'
var posPatternFindReg = regexp.MustCompile(" +pattern='([^']+)'")

// 取值表达式索引，格式为：patternIdx='n[:g]'，n和g默认值都是0
// n表示取匹配组中的第几个匹配值，默认0表示整个表达式匹配值
// g表示取第几组匹配值，默认0表示取所有匹配组
var posPatternIndexFindReg = regexp.MustCompile(" +patternIdx='((\\d+)([:]\\d+)?)'")
var posPatternIndexSplitReg = regexp.MustCompile("\\d+")

// 格式表达式，如日期格式，格式为：format='...'
var posFormatFindReg = regexp.MustCompile(" +format='([^']+)'")

// 时区表达式，默认时区为上海，格式为：timezone='...'
var posTimezoneFindReg = regexp.MustCompile(" +timezone='([^']+)'")

// 每个map表示一个item数据，key表示field index，value表示内容字符串，
// 由FieldUnmarshalTag负责解析value并生成元素对象，
type FieldValueMap map[int]reflect.Value

// 用于读取一个字段的内容字符串，实现方法由具体场景业务提供
type FieldTagReadValueFunc func(fieldTag *FieldTag, vars *Vars) (v string, err error)

type FieldTag struct {
	Id              int            //field index
	FieldName       string         //field name
	FieldType       reflect.Type   //
	Schema          string         // zip | xls | xlsx | xpath | csv
	Path            string         //position expression; e.g. sheet[x:x]/row[3:x]/col[3]
	PathFilled      string         //position expression filled vars; e.g. sheet[1]/row[2]/col[3]
	Pattern         *regexp.Regexp //find value
	PatternIdx      int            //find value, default 0
	PatternGroupIdx int            //find value, default 0
	Format          string         //format value, only for time
	Timezone        *time.Location //default +8, only for time
	Vars            *VarMap        //var map in Path
}
type FieldTagMap map[int]*FieldTag

func (m *FieldTagMap) MergeVars() *VarMap {
	if m == nil {
		return nil
	}
	var varMaps VarMaps
	for _, t := range *m {
		varMaps = append(varMaps, t.Vars)
	}
	return varMaps.Merge()
}
func (m *FieldTagMap) SetValues(vs *Vars) *FieldTagMap {
	if m == nil {
		return nil
	}
	for _, tag := range *m {
		tag.SetValues(vs)
	}
	return m
}
func (m *FieldTagMap) Filter(f func(tag *FieldTag) bool) *FieldTagMap {
	if m == nil {
		return nil
	}
	var re = make(FieldTagMap)
	for idx, t := range *m {
		if f(t) {
			re[idx] = t
		}
	}
	return &re
}

// 获取所有字段的FieldTag map，允许定义多个path schema
// T，目标对象类型
// tag, unmarshal tag is 'xm'
// varPattern, 定义path中的变量格式，表达式格式为：LTag(vName)[(vMin):(vMax)]RTag，例如 \{(\w+)\[(\d+):(\d+)(:(\d+))?\]\}
// 需确保匹配项长度为4，0可用于整体替换，1表示变量名，2表示变量下限，3表示变量上限
func GetFieldTags(T reflect.Type, tag string, varPattern *VarPattern) (*FieldTagMap, error) {
	var result = make(FieldTagMap)
	if varPattern == nil {
		varPattern = &DefaultVarPattern
	}
	var t = T
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fieldNum := t.NumField()
	//var lastProtocol = ""
	for i := 0; i < fieldNum; i++ {
		f := t.Field(i)
		posConfig, ok := f.Tag.Lookup(tag)
		if !ok {
			continue
		}
		path := strings.TrimSpace(posConfig)
		finds := posSchemaFindReg.FindStringSubmatch(path)
		if len(finds) != 3 {
			return nil, errors.New("bad field tag " + path)
		}
		schema := finds[1]
		//if lastProtocol == "" {
		//	lastProtocol = schema
		//} else if lastProtocol != schema {
		//	return nil, errors.New(fmt.Sprintf("multiple schemas not supported, %s, %s ", lastProtocol, schema))
		//}

		path = finds[2]
		path, patternStr := splitPosConfig(path, posPatternFindReg, "")
		path, patternIdxStr := splitPosConfig(path, posPatternIndexFindReg, "")
		path, formatStr := splitPosConfig(path, posFormatFindReg, "")
		path, timezoneStr := splitPosConfig(path, posTimezoneFindReg, "Asia/Shanghai")
		location, err := time.LoadLocation(timezoneStr)
		if err != nil {
			return nil, err
		}
		var pattern *regexp.Regexp = nil
		if patternStr != "" {
			pattern = regexp.MustCompile(patternStr)
		}

		patternIdxMatchStr := posPatternIndexSplitReg.FindAllString(patternIdxStr, -1)
		var patternIdx, patternGroupIdx int
		if len(patternIdxMatchStr) > 0 {
			patternIdx, _ = strconv.Atoi(patternIdxMatchStr[0])
		}
		if len(patternIdxMatchStr) > 1 {
			patternGroupIdx, _ = strconv.Atoi(patternIdxMatchStr[1])
		}

		result[i] =
			&FieldTag{
				Id:              i,
				FieldName:       f.Name,
				FieldType:       f.Type,
				Schema:          schema,
				Path:            path,
				Pattern:         pattern,
				PatternIdx:      patternIdx,
				PatternGroupIdx: patternGroupIdx,
				Format:          formatStr,
				Timezone:        location,
				Vars:            varPattern.Match(path),
			}
	}
	if len(result) == 0 {
		return nil, errors.New(tag + " tag not found")
	}
	return &result, nil
}

// 将变量填充到path
func (tag *FieldTag) SetValues(vs *Vars) *FieldTag {
	if tag == nil {
		return nil
	}
	if vs != nil && len(*vs) > 0 {
		var pathFilled = tag.Path
		for _, v := range *vs {
			if tv, ok := (*tag.Vars)[v.Name]; ok {
				if tv.Min <= v.Val && v.Val <= tv.Max && (v.Val-tv.Min)%tv.Step == 0 {
					tv.Val = v.Val
					pathFilled = strings.ReplaceAll(pathFilled, tv.Match, strconv.Itoa(tv.Val))
				} else {
					pathFilled = ""
				}
			}
		}
		tag.PathFilled = pathFilled
	}
	return tag
}

// 将读取的字符串内容，转换为字段类型的值
func (tag *FieldTag) Parse(str string) (reflect.Value, error) {
	if tag == nil {
		return reflect.Value{}, errors.New("FieldTag is nil")
	}
	if tag.Pattern != nil {
		allStringSubmatch := tag.Pattern.FindAllStringSubmatch(str, -1)
		str = ""
		if len(allStringSubmatch) > 0 && len(allStringSubmatch) > tag.PatternGroupIdx {
			if tag.PatternGroupIdx > 0 {
				if len(allStringSubmatch[tag.PatternGroupIdx]) > tag.PatternIdx && allStringSubmatch[tag.PatternGroupIdx][0] != "" {
					str = allStringSubmatch[tag.PatternGroupIdx][tag.PatternIdx]
				}
			} else {
				for t := range allStringSubmatch {
					if len(allStringSubmatch[t]) > tag.PatternIdx && allStringSubmatch[t][0] != "" {
						str += allStringSubmatch[t][tag.PatternIdx]
					}
				}
			}
		}
	}
	if str == "" && tag.FieldType.Kind() != reflect.String {
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
	return reflect.Value{}, errors.New("filed type not support " + tag.FieldType.Name())
}

func splitPosConfig(str string, pattern *regexp.Regexp, defaultV string) (string, string) {
	if pattern == nil {
		return str, defaultV
	}
	subMatches := pattern.FindStringSubmatch(str)
	if len(subMatches) > 1 {
		return strings.Replace(str, subMatches[0], "", 1), subMatches[1]
	} else {
		return str, defaultV
	}
}
