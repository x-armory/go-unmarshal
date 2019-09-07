package base

// 用于处理Item，可控制是否继续处理流程
type ItemFilter func(item interface{}, vars *Vars) (flow FlowControl, deep int)
type ItemFilters []ItemFilter

func (filters ItemFilters) Filter(item interface{}, vars *Vars) (flow FlowControl, deep int) {
	for _, filter := range filters {
		flow, deep := filter(item, vars)
		if flow != Forward {
			return flow, deep
		}
	}
	return Forward, deep
}

func ExitIfFieldValueStringMapIsEmpty(item interface{}, vars *Vars) (flow FlowControl, deep int) {
	model, ok := item.(*FieldValueMap)
	if !ok {
		return Forward, 0
	}
	if len(*model) == 0 {
		println("[WARN] ExitIfFieldValueStringMapIsEmpty")
		return Break, 0
	}
	return Forward, 0
}
func ExitIfItemIsNil(item interface{}, vars *Vars) (flow FlowControl, deep int) {
	if item == nil {
		return Break, 0
	}
	return Forward, 0
}
