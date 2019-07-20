package mhttp

// Encoder describes a type capable of encoding a value into a []byte.
type Encoder func(interface{}) ([]byte, error)
