package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

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

// MustNewRequest fails with ginkgo.Fail if making request fails
func MustNewRequest(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return req
}

// MustURLParse fails with ginkgo.Fail if making request fails
func MustURLParse(urlStr string) *url.URL {
	url, err := url.Parse(urlStr)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return url
}
