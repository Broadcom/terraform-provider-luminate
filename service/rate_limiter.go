// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"golang.org/x/time/rate"
	"net/http"
)

type RateLimitTransport struct {
	limiter      *rate.Limiter
	roundTripper http.RoundTripper
}

func NewRateLimitTransport(requestsPerSecond float64, burst int, roundTripper http.RoundTripper) http.RoundTripper {
	return &RateLimitTransport{
		limiter:      rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
		roundTripper: roundTripper,
	}
}

func (t *RateLimitTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.limiter.Wait(r.Context())
	return t.roundTripper.RoundTrip(r)
}
