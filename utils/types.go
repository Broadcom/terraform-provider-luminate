package utils

import sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"

// List of WebRDP ApplicationSubType
const (
	RDP_BROWSER_SINGLE_MACHINE_ApplicationSubType    sdk.ApplicationSubType = "RDP_BROWSER_SINGLE_MACHINE"
	RDP_BROWSER_MULTIPLE_MACHINES_ApplicationSubType sdk.ApplicationSubType = "RDP_BROWSER_MULTIPLE_MACHINES"

	RDP_NATIVE_AccessType  string = "RDP_NATIVE"
	RDP_BROWSER_AccessType string = "RDP_BROWSER"
)
