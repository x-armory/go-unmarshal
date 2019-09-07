package base

import "io"

type Unmarshaler interface {
	Unmarshal(r io.Reader, data interface{}) error
}
