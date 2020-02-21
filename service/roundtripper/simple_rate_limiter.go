package roundtripper

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

type SimpleRateLimitTransport struct {
	roundTripper http.RoundTripper
	sleepBetweenRequest time.Duration
}

func NewSimpleRateLimitTransport(requestsPerSecond float64, roundTripper http.RoundTripper) http.RoundTripper {

	sleepBetweenRequest:= math.Round(1000 / requestsPerSecond)

	return &SimpleRateLimitTransport{
		roundTripper: roundTripper,
		sleepBetweenRequest: time.Duration(sleepBetweenRequest),
	}
}

func (t *SimpleRateLimitTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Println("Send request on: " + time.Now().String())
	time.Sleep(t.sleepBetweenRequest *  time.Millisecond)
	return t.roundTripper.RoundTrip(r)
}
