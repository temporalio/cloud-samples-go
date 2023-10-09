package api

import _ "embed"

const (
	TemporalCloudAPIVersionHeader = "temporal-cloud-api-version"
)

var (
	//go:embed version
	TemporalCloudAPIVersion string
)
