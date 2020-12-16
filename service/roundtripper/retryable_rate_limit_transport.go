package roundtripper

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type RetryableRateLimitTransport struct {
	roundTripper http.RoundTripper
	sleepBetweenRetries time.Duration
	retrySleepJitter int
}

func NewRetryableRateLimitTransport(millsBetweenRetries int, retrySleepJitter int, roundTripper http.RoundTripper) http.RoundTripper {
	rand.Seed(time.Now().UnixNano())

	return &RetryableRateLimitTransport{
		roundTripper: roundTripper,
		sleepBetweenRetries: time.Duration(millsBetweenRetries) * time.Millisecond,
		retrySleepJitter: retrySleepJitter,
	}
}

func (t *RetryableRateLimitTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	// Duplicate the request
	var body bytes.Buffer
	if r.ContentLength > 0 {
		body.ReadFrom(r.Body)
	}

	r2 := r.Clone(context.Background())
	r2.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))

	r.Body = ioutil.NopCloser(&body)

	response, err := t.roundTripper.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == 429 {
		log.Printf("[DEBUG] Retry request due to rate limit error.")

		// Jitter sleep time to spread the retries more evenly
		sleep := t.sleepBetweenRetries + (time.Duration(rand.Intn(t.retrySleepJitter)) * time.Millisecond) - (time.Duration(t.retrySleepJitter / 2) * time.Millisecond)
		time.Sleep(sleep)

		return t.RoundTrip(r2)
	}
	return response, err
}
