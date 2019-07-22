// Package httpm provides functions to produce and consume http messages.
//
//	It's main goal is to reduce the boilerplate due to:
//	* writing request and reading response on client side
//	* reading request and writing response on server side
//
//	The ComposeXXX functions return the possible error that occures during composition,
//	leaving a very clean describtion of what to do when composing or consuming a message.
//
//	// Classic way ...
//	raw, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		...
//	}
//	if err := json.Unmarshal(raw, &target); err != nil {
//		...
//	}
//	if err := check(target); err != nil {
//		...
//	}
//	if err := parseParams(r.URL, &params); err != nil {
//		...
//	}
//
//	// With httpm ...
//	parse := httpm.ComposeRequest(
//		httpm.ReadRequestBody(httpm.DecodeAndCheck(json.Unmarshal, check),
//		httpm.ReadRequestParams(parseParams),
//	)
//	r, err := parse(r)
//	if err != nil {
//		...
//	}
package httpm
