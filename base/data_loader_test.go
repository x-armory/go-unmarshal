package base

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

type GetDataLoaderModel struct {
	Name string     `xm:"zip:{Name}"`
	F1   string     `xm:"excel:sheet[5:7]/row[3:x]/col[3]"`
	F2   time.Time  `xm:"excel:sheet[6]/row[:]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"`
	F3   *time.Time `xm:"excel:sheet[6]/row[:]/col[3] pattern='\\d{4}-\\d{2}-\\d{2}' format='2006-01-02' timezone='Asia/Shanghai'"`
	F4   int        `xm:"excel:sheet[6]/row[:]/col[3] pattern='\\d+'"`
}

func TestGetDataLoader(t *testing.T) {
	var data []*GetDataLoaderModel
	dataLoader := DataLoader{
		Data:      &data,
		WriteData: false,
		VarOrder:  []string{"sheet", "row", "col"},
		ReadValueFunc: map[string]FieldTagReadValueFunc{
			"excel": func(fieldTag *FieldTag, vars *Vars) (v string, err error) {
				return fmt.Sprintf("%s = %d", fieldTag.Path+" "+vars.ToString(), rand.Int63()), nil
			},
			"zip": func(fieldTag *FieldTag, vars *Vars) (v string, err error) {
				return fmt.Sprintf("%s = %d", fieldTag.Path+" "+vars.ToString(), rand.Int63()), nil
			},
		},
		ItemFilters: []ItemFilter{
			func(item interface{}, vars *Vars) (flow FlowControl, deep int) {
				model, ok := item.(*GetDataLoaderModel)
				if !ok {
					return Forward, 1
				}
				println("get an item:")
				// process data item
				// ignore row 5
				if (*vars)[1].Val == 5 {
					return Continue, 1
				}
				//
				if (*vars)[1].Val > 10 {
					return Continue, 0
				}
				//
				if (*vars)[0].Val > 6 && (*vars)[1].Val == 3 {
					return Break, 0
				}

				println(fmt.Sprintf("%+v", model))
				return Forward, 1
			},
		},
	}
	e := dataLoader.Load()
	assert.NoError(t, e)
}
