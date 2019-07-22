# TODO:

## httpm

- README.md redaction
- Package documentation

- ~Request (Q)~
- ~Response (R)~
- ResponseWriter (W)
	- Tests
	- Examples
	- ResponseWriter with Code implementation, and associated middleware
- Decoder
	- ExtendDecoderFn(decoder, DecoderExtensionFn) Decoder
	- Validator.Validate as DecoderExtensionFn
	- Schema.Parse as DecoderExtension

- Model binding middleware (low priority)
	- the parameter name
	- the conversion func from name to key -- mux.Vars or httprouter.P
	- the function returning object from key -- mapper.GetXXX
		- func(error) []byte -- default to http.StatusText
		- func(error) http status -- default to 404
	- the context getter/setter (iterable) -- CtxSetXXX


## mcontext

- README.md redaction
- Package documentation

- Code recuperation from bitbucket
- Tests
- Examples

## msqlx

- README.md redaction
- Package documentation

- Code recuperation from bitbucket
- Tests
- Examples

## contributing

- define rules for contribution
- setup CI (travis ?)
