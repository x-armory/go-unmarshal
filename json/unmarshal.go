package json

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Unmarshaler struct {
	IgnoreError bool
}

func (m *Unmarshaler) Unmarshal(r io.Reader, data interface{}) error {
	var rt = *m
	bts, e := ioutil.ReadAll(r)
	if e != nil {
		return e
	}
	if e := json.Unmarshal(bts, data); e != nil {
		if !rt.IgnoreError {
			return e
		}
		println("[WARN]", e.Error())
	}
	return nil
}
