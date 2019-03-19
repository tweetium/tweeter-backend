package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/onsi/ginkgo"
)

// MustJSONMarshal fails with ginkgo.Fail if marshalling fails
// Returns *bytes.Buffer to be used in a request
func MustJSONMarshal(v interface{}) *bytes.Buffer {
	buf, err := json.Marshal(v)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return bytes.NewBuffer(buf)
}

// MustJSONUnmarshalStrict disallows unknown fields on the json provided and
// fails if the unmarshals errors
func MustJSONUnmarshalStrict(b []byte, v interface{}) {
	// This check is not done when using json.Decoder directly (instead of json.Unmarshal)
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		ginkgo.Fail(fmt.Sprintf("Can't unmarshal non-pointer type: '%s'", reflect.TypeOf(v).Name()))
	}

	dec := json.NewDecoder(bytes.NewBuffer(b))
	dec.DisallowUnknownFields()
	err := dec.Decode(&v)
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Failed to unmarshal as type: '%s', raw: %s, err: %s", reflect.TypeOf(v).Name(), string(b), err))
	}
}

// MustNewRequest fails with ginkgo.Fail if making request fails
func MustNewRequest(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return req
}
