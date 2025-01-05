// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func expandStringList(configured *schema.Set) []*string {
	configuredList := configured.List()
	vs := make([]*string, 0, len(configuredList))
	for _, v := range configuredList {
		val, ok := v.(string)
		if ok && val != "" {
			str := v.(string)
			vs = append(vs, &str)
		}
	}
	return vs
}
