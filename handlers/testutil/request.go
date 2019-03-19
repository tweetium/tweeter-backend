package testutil

import (
	"fmt"
	"io"
	"strings"
	"tweeter/testutil"
)

// RequestArgs is the easily modified version of an http.Request that can
// be used with other testutil functions to send requests
type RequestArgs struct {
	Method   string
	Endpoint string
	JSONBody map[string]interface{}
	RawBody  *string
}

// GetBody returns the body of the request and validates that multiple valid
// bodies are not supplied
func (r RequestArgs) GetBody() (body io.Reader, err error) {
	countValid := 0
	if r.JSONBody != nil {
		body = testutil.MustJSONMarshal(r.JSONBody)
		countValid++
	}
	if r.RawBody != nil {
		body = strings.NewReader(*r.RawBody)
		countValid++
	}
	if countValid > 1 {
		err = fmt.Errorf("RequestArgs has multiple valid bodies, request: %+v", r)
		return
	}

	if countValid == 0 {
		err = fmt.Errorf("RequestArgs has no valid bodies, request: %+v", r)
		return
	}

	return
}
