package base

import (
	"fmt"
	"regexp"
	"strconv"
)

// FieldName：变量名，出现重名，会被覆盖；
// Match：匹配到的表达式，用于整体替换；
// Min：变量下线；
// Max：变量上限；
// Val：自动生成的遍历值
// Step：增长步长，默认为1
type Var struct {
	Name  string
	Match string
	Min   int
	Max   int
	Val   int
	Step  int
}

func (v *Var) IsValid() bool {
	if v == nil {
		return false
	}
	return v.Val <= v.Max
}
func (v *Var) Reset() {
	v.Val = v.Min
}
func (v *Var) Next() {
	if v.Step == 0 {
		v.Step = 1
	}
	v.Val += v.Step
}

// set          vars := [ v0,        v1,        v2,        v2                     ]
// then                   v0.deep=0, v1.deep=1, v2.deep=2, v2.deep=3
//
// break    (deep=0) -> { v0.Max(),  v1.Max(),  v2.Max(),  v3.Max(),  vars.Next() }
// continue (deep=0) -> {            v1.Max(),  v2.Max(),  v3.Max(),  vars.Next() }
//
// break    (deep=1) -> {            v1.Max(),  v2.Max(),  v3.Max(),  vars.Next() }
// continue (deep=1) -> {                       v2.Max(),  v3.Max(),  vars.Next() }
//
// break    (deep=2) -> {                       v2.Max(),  v3.Max(),  vars.Next() }
// continue (deep=2) -> {                                  v3.Max(),  vars.Next() }
//
// break    (deep=3) -> {                                  v3.Max(),  vars.Next() }
// continue (deep=3) -> {                                             vars.Next() }
type Vars []*Var

func (vs *Vars) IsValid() bool {
	if vs == nil || len(*vs) == 0 {
		return false
	}
	return (*vs)[0].IsValid()
}
func (vs *Vars) Reset() {
	if vs == nil || len(*vs) == 0 {
		return
	}
	for _, v := range *vs {
		v.Reset()
	}
}
func (vs *Vars) Next() {
	if vs == nil || !vs.IsValid() {
		return
	}
	for i := len(*vs) - 1; i >= 0; i-- {
		v := (*vs)[i]
		v.Next()
		if v.Val > v.Max {
			if i > 0 {
				v.Val = v.Min
			} else {
				return
			}
		} else {
			return
		}
	}
}
func (vs *Vars) setMax(deep int) {
	for i := deep; i <= len(*vs)-1; i++ {
		(*vs)[i].Val = (*vs)[i].Max
	}
}
func (vs *Vars) Break(deep int) {
	vs.setMax(deep)
	vs.Next()
}
func (vs *Vars) Continue(deep int) error {
	vs.setMax(deep + 1)
	vs.Next()
	return nil
}

func (vs *Vars) ToString() string {
	var re = "Vars:"
	if vs == nil {
		return re
	}
	for _, v := range *vs {
		re += fmt.Sprintf(" %s=%d(%d-%d)", v.Name, v.Val, v.Min, v.Max)
	}
	return re
}

type VarMap map[string]*Var

// 将VarMap转化为数组
// 如果front不为空，则将front先加入数组
func (ms *VarMap) List(front ...string) *Vars {
	if ms == nil {
		return nil
	}
	var re = make(Vars, 0)
	var addedName = make(map[string]bool)
	for _, n := range front {
		if v, ok := (*ms)[n]; ok {
			re = append(re, v)
			addedName[n] = true
		}
	}
	for _, m := range *ms {
		if _, ok := addedName[m.Name]; !ok {
			re = append(re, m)
		}
	}
	return &re
}

type VarMaps []*VarMap

// 用于合并相同变量的上下限，取得循环范围
func (ms *VarMaps) Merge() *VarMap {
	if ms == nil {
		return nil
	}
	var re = make(VarMap)
	for _, m := range *ms {
		for _, v := range *m {
			if reV, ok := re[v.Name]; !ok {
				re[v.Name] = &Var{
					Name: v.Name,
					Min:  v.Min,
					Max:  v.Max,
					Step: v.Step,
				}
			} else {
				if reV.Min > v.Min {
					reV.Min = v.Min
				}
				if reV.Max < v.Max {
					reV.Max = v.Max
				}
				if reV.Step != 0 && v.Step != 0 && reV.Step != v.Step {
					reV.Step = 1
				}
			}
		}
	}
	return &re
}

// Pattern：变量正则表达式，用于匹配变量定义
// ParseFunc：将匹配的字符串数组，转化为Var对象
type VarPattern struct {
	Pattern   *regexp.Regexp
	ParseFunc VarParseFunc
}
type VarParseFunc func(v []string) *Var

// 从path表达式中获取变量map
func (p VarPattern) Match(s string) *VarMap {
	var re = make(VarMap)
	matches := p.Pattern.FindAllStringSubmatch(s, -1)
	for _, m := range matches {
		if v := p.ParseFunc(m); v != nil {
			re[v.Name] = v
		}
	}
	return &re
}

// 默认的变量匹配表达式，格式为：VarName[min:max:step]，min和max可以为空，默认值为0，999999
var DefaultVarPattern = VarPattern{
	Pattern: regexp.MustCompile(`(\w+)\[(\d+)?:(\d+)?(:\d+)?\]`),
	ParseFunc: func(v []string) *Var {
		if len(v) != 5 {
			return nil
		}
		var min = 0
		var max = 999999
		var step = 1
		if v[2] != "" {
			min, _ = strconv.Atoi(v[2])
		}
		if v[3] != "" {
			max, _ = strconv.Atoi(v[3])
		}
		if v[4] != "" {
			step, _ = strconv.Atoi(v[4][1:])
		}
		if step == 0 {
			step = 1
		}
		if max < min {
			min, max = max, min
		}
		return &Var{
			Name:  v[1],
			Match: v[0],
			Min:   min,
			Max:   max,
			Step:  step,
		}
	},
}
