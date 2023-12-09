//go:build !release

package config

// Use these values sparingly. Anything added here should be also added to release.go.
var (
	studioProxy = "http://localhost:8001"
)
