package json

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Unmarshaler struct {
	IgnoreError bool
}

func (m *Unmarshaler) unmarshal(r io.Reader, data interface{}) error {
	bts, e := ioutil.ReadAll(r)
	if e != nil {
		return e
	}
	if e := json.Unmarshal(bts, data); e != nil {
		if !m.IgnoreError {
			return e
		}
		println("[WARN]", e.Error())
	}
	return nil
}
