package integration

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg"
	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/config"
	httpPresentation "gses2.app/api/pkg/presentation/http"
)

func TestThatRateRouteReturnsStatusOKAndFloat(t *testing.T) {
	config.LoadEnv()

	testServer := getTestServerWithRateRoute()
	defer testServer.Close()

	response, err := http.Get(testServer.URL)
	defer closeResponseBody(t, *response)
	rate := getResponseBodyContent(t, *response)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assertIsFloat(t, rate)
}

func TestThatGetExchangeRateReturnValuesDontFluctuateMuchOnSuccessiveCallsAfterOneSecond(t *testing.T) {
	config.LoadEnv()

	oneSecondDuration := time.Duration(1_000_000_000)

	testServer := getTestServerWithRateRoute()
	defer testServer.Close()

	response1, err1 := http.Get(testServer.URL)
	defer closeResponseBody(t, *response1)

	time.Sleep(oneSecondDuration)

	response2, err2 := http.Get(testServer.URL)
	defer closeResponseBody(t, *response2)

	rate1, _ := strconv.ParseFloat(getResponseBodyContent(t, *response1), 64)
	rate2, _ := strconv.ParseFloat(getResponseBodyContent(t, *response2), 64)
	delta := rate1 * 0.25 // delta is 25% of rate1

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.InDelta(t, rate1, rate2, delta)
}

func getResponseBodyContent(t *testing.T, response http.Response) string {
	t.Helper()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Error when reading response body")
	}

	return string(content)
}

func getTestServerWithRateRoute() *httptest.Server {
	genericExchangeRateService := pkg.GetGenericExchangeRateService()
	btcToUahService := application.NewBtcToUahServiceImpl(genericExchangeRateService)
	rateHandler := httpPresentation.RateRequestHandler{BtcToUahService: btcToUahService}
	testServer := httptest.NewServer(http.HandlerFunc(httpPresentation.GetHandlerFunction(rateHandler)))

	return testServer
}

func closeResponseBody(t *testing.T, response http.Response) {
	t.Helper()

	err := response.Body.Close()
	if err != nil {
		t.Error("Failure to close response body")
	}
}

func assertIsFloat(t *testing.T, str string) {
	t.Helper()

	_, err := strconv.ParseFloat(str, 64)
	if err != nil {
		t.Error(str, "is not float")
	}
}
