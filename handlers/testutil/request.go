package testutil

// RequestArgs is the easily modified version of an http.Request that can
// be used with other testutil functions to send requests
type RequestArgs struct {
	Method   string
	Endpoint string
	JSONBody map[string]interface{}
}
