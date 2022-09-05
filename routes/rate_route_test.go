package routes

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/config"
)

func TestThatRateRouteReturnsStatusOKAndFloat(t *testing.T) {
	config.LoadEnv()
	testServer := httptest.NewServer(http.HandlerFunc(rateRoute))
	defer testServer.Close()

	response, err := http.Get(testServer.URL)
	rate := getResponseBodyContent(t, *response)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assertIsFloat(t, rate)
}

func TestThatGetExchangeRateReturnValuesDontFluctuateMuchOnSuccessiveCallsAfterOneSecond(t *testing.T) {
	config.LoadEnv()
	oneSecondDuration := time.Duration(1_000_000_000)
	testServer := httptest.NewServer(http.HandlerFunc(rateRoute))
	defer testServer.Close()

	response1, err1 := http.Get(testServer.URL)
	time.Sleep(oneSecondDuration)
	response2, err2 := http.Get(testServer.URL)

	rate1, _ := strconv.ParseFloat(getResponseBodyContent(t, *response1), 64)
	rate2, _ := strconv.ParseFloat(getResponseBodyContent(t, *response2), 64)
	delta := rate1 * 0.25 // delta is 25% of rate1

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.InDelta(t, rate1, rate2, delta)
}

func getResponseBodyContent(t *testing.T, response http.Response) string {
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Error when reading response body")
	}

	err = response.Body.Close()
	if err != nil {
		t.Error("Failure to close response body")
	}

	return string(content)
}

func assertIsFloat(t *testing.T, str string) {
	_, err := strconv.ParseFloat(str, 64)
	if err != nil {
		t.Error(str, "is not float")
	}
}
