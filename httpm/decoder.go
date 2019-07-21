package httpm

// Decoder describes a type capable of decoding a []byte.
type Decoder func([]byte, interface{}) error
