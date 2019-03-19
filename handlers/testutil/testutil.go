package testutil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/onsi/ginkgo"

	"tweeter/handlers/responses"
	"tweeter/testutil"
)

// MustReadSuccessData fails with ginkgo.Fail if reading success request fails
func MustReadSuccessData(resp *http.Response, data interface{}) {
	var successResp responses.SuccessResponse
	successResp.Data = data

	mustJSONUnmarshalResponse(resp, &successResp)
}

// MustReadErrors fails with ginkgo.Fail if reading error request fails
func MustReadErrors(resp *http.Response) []responses.Error {
	var errorResp responses.ErrorResponse
	mustJSONUnmarshalResponse(resp, &errorResp)

	return errorResp.Errors
}

// MustSendRequest sends a correctly formatted request to the httptest.Server and fails if any error
func MustSendRequest(server *httptest.Server, request RequestArgs) *http.Response {
	url := server.URL + request.Endpoint
	body, err := request.GetBody()
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Request body is invalid, err: %s, request: %+v", err, request))
	}
	httpRequest := testutil.MustNewRequest(request.Method, url, body)

	resp, err := server.Client().Do(httpRequest)
	if err != nil {
		errString := err.Error()
		// if you try to send the same request twice, the body is consumed
		if strings.Contains(errString, "with Body length 0") {
			ginkgo.Fail(fmt.Sprintf("Failed to send request (looks like the request was sent twice) with err: %s", errString))
		} else {
			ginkgo.Fail(errString)
		}
	}

	return resp
}

func mustJSONUnmarshalResponse(resp *http.Response, v interface{}) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	testutil.MustJSONUnmarshalStrict(bodyBytes, v)
}
