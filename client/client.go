// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package client

type SecureCloudClient struct {
	ApiEndpoint      string
	ApiClientID      string
	APIClientSeceret string
}

func NewClient(apiEndpoint string, apiClientID string, apiClientSecret string) (SecureCloudClient, error) {
	cli := SecureCloudClient{
		ApiEndpoint:      apiEndpoint,
		ApiClientID:      apiClientID,
		APIClientSeceret: apiClientSecret,
	}
	return cli, nil
}
