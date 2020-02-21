package roundtripper

import (
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type RetryableRateLimitTransport struct {
	roundTripper http.RoundTripper
	sleepBetweenRetries time.Duration
}

func NewRetryableRateLimitTransport(millsBetweenRetries int, roundTripper http.RoundTripper) http.RoundTripper {

	return &RetryableRateLimitTransport{
		roundTripper: roundTripper,
		sleepBetweenRetries: time.Duration(millsBetweenRetries) * time.Millisecond,
	}
}

func (t *RetryableRateLimitTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	response, err := t.roundTripper.RoundTrip(r)

	if err != nil {
		switch err.(type) {
		case *oauth2.RetrieveError:
			// in case of rate-limit error, retry after sleepBetweenRetries duration
			if err.(*oauth2.RetrieveError).Response.StatusCode == 429 {
				fmt.Println("Retry request due to rate limit error.")
				time.Sleep(t.sleepBetweenRetries)
				return t.roundTripper.RoundTrip(r)
			}
		default:
			return nil, err
		}
	}

	return response, err
}
