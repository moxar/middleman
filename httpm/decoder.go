package httpm

import (
	"encoding"
	"fmt"
)

// Decoder describes a type capable of decoding a []byte.
type Decoder func([]byte, interface{}) error

var (
	_ Decoder = DecodeText
)

// DecodeText decodes the raw data into the input.
// Input must be of type *[]byte, *string, or encoding.TextUnmarshaler
func DecodeText(raw []byte, v interface{}) error {
	switch t := v.(type) {
	case *string:
		*t = string(raw)
		return nil
	case *[]byte:
		*t = raw
		return nil
	case encoding.TextUnmarshaler:
		return t.UnmarshalText(raw)
	default:
		return fmt.Errorf("cannot decode %T into %T: %T does not implement encoding.TextUnmarshaller", []byte{}, v, v)
	}
}
