package xlsx

import "testing"

type Ref struct {
	Obj
}
type Obj struct {
	A string
}

func (r *Ref) Set(s string) {
	r.A = s
}
func TestObjectRef(t *testing.T) {
	ref := &Ref{Obj{"Init"}}
	ref1 := &Ref{ref.Obj}
	ref2 := &Ref{ref.Obj}
	ref1.Set("1")
	ref2.Set("2")
	println(ref.A, ref1.A, ref2.A)
}
