package csv

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"github.com/x-armory/go-unmarshal/base"
	"os"
	"regexp"
	"strings"
	"testing"
)

type MyDto struct {
	Seq        int    `xm:"csv://row[r[0:]]/col[0] pattern='\\d+'"`
	MemberName string `xm:"csv://row[r[0:]]/col[2]"`
	Vol        int    `xm:"csv://row[r[0:]]/col[3] pattern='[-]?\\d+'"`
	Add        int    `xm:"csv://row[r[0:]]/col[5] pattern='[-]?\\d+'"`
}

func TestGetVar_(t *testing.T) {
	r := regexp.MustCompile(`//(row\[)?(\d+)(\])?/(col\[)?(\d+)(\])?`).MatchString("//row[2]/col[1]")
	println(r)
}

func TestUnmarshaler_Unmarshal(t *testing.T) {
	file, e := os.Open("/private/tmp/20191022_DCE_DPL/20191022_a2001_成交量_买持仓_卖持仓排名.txt")
	assert.NoError(t, e)

	unmarshaler := Unmarshaler{
		DataLoader: base.DataLoader{
			WriteData:       true,
			ExitNoDataTimes: 10,
			VarOrder:        []string{"r", "col"},
			ItemFilters: []base.ItemFilter{
				func(item interface{}, vars *base.Vars) (flow base.FlowControl, deep int) {
					data, ok := item.(*MyDto)
					if !ok {
						return base.Forward, 0
					}
					if data.Seq == 0 || data.Vol == 0 {
						return base.Continue, 0
					}
					return base.Forward, 0
				},
			},
		},
	}

	data := []*MyDto{}
	e = unmarshaler.Unmarshal(file, &data)
	assert.NoError(t, e)
	println(len(data))
}

func TestReader(t *testing.T) {
	file, e := os.Open("/private/tmp/20191022_DCE_DPL/20191022_a2001_成交量_买持仓_卖持仓排名.txt")
	assert.NoError(t, e)
	reader := bufio.NewReader(file)

	for true {
		line, _, e := reader.ReadLine()
		if e != nil {
			println(e.Error())
			break
		}
		content := string(line)
		split := strings.Split(content, "\t")
		if len(split) != 8 {
			continue
		}
		print(len(split), " ")
		for i := range split {
			print("\t(", i+1, ")", split[i])
		}
		println()
	}
}
