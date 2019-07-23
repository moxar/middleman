# TODO:

## httpm

- Response (R)
	- Tests
	- Examples
- ResponseWriter (W)
	- Tests
	- Examples

- Model binding middleware (low priority)
	- the parameter name
	- the conversion func from name to key -- mux.Vars or httprouter.P
	- the function returning object from key -- mapper.GetXXX
		- func(error) []byte -- default to http.StatusText
		- func(error) http status -- default to 404
	- the context getter/setter (iterable) -- CtxSetXXX

## mcontext

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
