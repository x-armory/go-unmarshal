package base

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

func TestGetVar1(t *testing.T) {
	path := "//*[@id=\"maincontent\"]/div[4]/table/tbody/tr[{row[2:3]}]/td[{col[2:3]}]"
	varPattern := VarPattern{
		Pattern: regexp.MustCompile(`\{(\w+)\[(\d+):(\d+)\]\}`),
		ParseFunc: func(v []string) *Var {
			if len(v) != 4 {
				return nil
			}
			min, _ := strconv.Atoi(v[2])
			max, _ := strconv.Atoi(v[3])
			return &Var{
				Name:  v[1],
				Match: v[0],
				Min:   min,
				Max:   max,
			}
		},
	}

	vars := varPattern.Match(path)
	for _, v := range *vars {
		fmt.Printf("%+v\n", v)
	}
}

func TestDefaultVarPattern(t *testing.T) {
	path := "sheet[]/row[3:]/col[:3]/item[0]"
	vars := DefaultVarPattern.Match(path)
	for _, v := range *vars {
		fmt.Printf("%+v\n", v)
	}
}

func TestDefaultVarPattern2(t *testing.T) {
	path := "//*[@id=\"maincontent\"]/div[4]/table/tbody/tr[row[2:3:2]]/td[col[2:3]]"
	vars := DefaultVarPattern.Match(path)
	for name, v := range *vars {
		fmt.Printf("%s -> %+v\n", name, *v)
	}
}

func TestVarMaps_Merge(t *testing.T) {
	var vms = make(VarMaps, 0)
	vms = append(vms, DefaultVarPattern.Match("//*[@id=\"maincontent\"]/div[4]/table/tbody/tr[row[2:3:2]]/td[col[2:3]]"))
	vms = append(vms, DefaultVarPattern.Match("//*[@id=\"maincontent\"]/div[4]/table/tbody/tr[row[2:40:2]]/td[col[20:1]]"))
	vms = append(vms, DefaultVarPattern.Match("//*[@id=\"maincontent\"]/div[4]/table/tbody/tr[row[1:3:2]]/td[col[7:30]]"))
	vars := vms.Merge()
	for name, v := range *vars {
		fmt.Printf("%s -> %+v\n", name, *v)
	}
}

func TestVar_Next(t *testing.T) {
	v := &Var{
		Name: "v",
		Min:  4,
		Max:  8,
	}
	for v.Reset(); v.IsValid(); v.Next() {
		fmt.Printf("%+v\n", v)
	}
}

func TestVars_Next(t *testing.T) {
	var vms = make(VarMaps, 0)
	vms = append(vms, DefaultVarPattern.Match("//*[@id=\"maincontent\"]/div[page[:1]]/table/tbody/tr[row[1:3]]/td[col[2:3]]"))
	vms = append(vms, DefaultVarPattern.Match("//*[@id=\"maincontent\"]/div[page[3:5]]/table/tbody/tr[row[2:5]]/td[col[8:1]]"))
	vms = append(vms, DefaultVarPattern.Match("//*[@id=\"maincontent\"]/div[4]/table/tbody/tr[row[1:3]]/td[col[7:5]]"))
	vars := vms.Merge().List("page", "row", "col")

	for vars.Reset(); vars.IsValid(); vars.Next() {
		for _, v := range *vars {
			print("\t", v.Val)
		}
		println()
	}
}

func TestVarMap_List(t *testing.T) {
	var vms = make(VarMaps, 0)
	vms = append(vms, DefaultVarPattern.Match("//*[@id=\"maincontent\"]/div[page[:1]]/table/tbody/tr[row[1:3]]/td[col[2:3]]"))
	vms = append(vms, DefaultVarPattern.Match("//*[@id=\"maincontent\"]/div[page[3:5]]/table/tbody/tr[row[2:5]]/td[col[8:1]]"))
	vms = append(vms, DefaultVarPattern.Match("//*[@id=\"maincontent\"]/div[4]/table/tbody/tr[row[1:3]]/td[col[7:5]]"))
	vars := vms.Merge().List("page", "row", "col")
	bts, _ := json.MarshalIndent(vars, "", "  ")
	println(string(bts))
}
