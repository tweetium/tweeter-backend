package testutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tweeter/handlers/responses"

	"github.com/onsi/ginkgo"
)

// MustReadSuccessData fails with ginkgo.Fail if reading success request fails
func MustReadSuccessData(resp *http.Response, data interface{}) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	var successResp responses.SuccessResponse
	successResp.Data = data
	err = json.Unmarshal(bodyBytes, &successResp)
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Failed to unmarshal as responses.SuccessResponse, raw: %s, err: %s", string(bodyBytes), err))
	}
}
