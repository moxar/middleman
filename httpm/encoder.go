package httpm

import (
	"encoding"
	"fmt"
)

var (
	_ Encoder = EncodeText
)

// Encoder describes a type capable of encoding a value into a []byte.
type Encoder func(interface{}) ([]byte, error)

// EncodeText encodes string, []byte and text.Marshaller implementations.
func EncodeText(v interface{}) ([]byte, error) {
	switch t := v.(type) {
	case string:
		return []byte(t), nil
	case []byte:
		return t, nil
	case encoding.TextMarshaler:
		return t.MarshalText()
	default:
		return nil, fmt.Errorf("cannot encode %T into %T: %T does not implement encoding.TextMarshaller", v, []byte{}, v)
	}
}
