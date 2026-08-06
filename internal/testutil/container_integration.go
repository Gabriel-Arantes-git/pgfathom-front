//go:build integration

// Package testutil, integration half.
//
// Everything behind this build tag needs Docker. Keeping it tagged is what lets
// `go test ./...` stay fast, offline and dependency-light: someone who wants to
// contribute a naming profile for their language should be able to clone, run
// the suite and send a patch without installing a container runtime.
//
// Phase 2 fills this in with a testcontainers-backed PostgreSQL fixture. The
// file exists now so the suite does not have to be reorganised then, and so the
// tag is proven to work from the start.
package testutil

// IntegrationEnabled reports whether the integration build tag is active. It
// exists so the tag itself is covered by a test rather than assumed.
const IntegrationEnabled = true
