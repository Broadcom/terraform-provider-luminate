// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package test_utils

import "math/rand"

func GetRandomNumber() int {
	return GetRandomNumberWithBase(10000)
}

func GetRandomNumberWithBase(n int) int {
	return n + rand.Intn(n)
}
