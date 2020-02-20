package service

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

type CustomRateLimitTransport struct {
	roundTripper http.RoundTripper
	sleepBetweenRequest time.Duration
}

func NewCustomRateLimitTransport(requestsPerSecond float64, roundTripper http.RoundTripper) http.RoundTripper {

	sleepBetweenRequest:= math.Round(1000 / requestsPerSecond)

	return &CustomRateLimitTransport{
		roundTripper: roundTripper,
		sleepBetweenRequest: time.Duration(sleepBetweenRequest),
	}
}

func (t *CustomRateLimitTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Println("Send request on: " + time.Now().String())
	time.Sleep(t.sleepBetweenRequest *  time.Millisecond)
	return t.roundTripper.RoundTrip(r)
}
