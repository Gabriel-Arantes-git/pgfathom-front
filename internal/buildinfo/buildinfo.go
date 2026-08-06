// Package buildinfo carries identity stamped in at link time.
package buildinfo

import "runtime/debug"

// Injected via -ldflags at build time. See the Makefile.
var (
	Version = ""
	Commit  = ""
	Date    = ""
)

const unknown = "unknown"

// Resolve returns version, commit and build date, falling back to what the Go
// toolchain recorded when built without ldflags — the `go install` case.
func Resolve() (version, commit, date string) {
	version, commit, date = Version, Commit, Date

	info, ok := debug.ReadBuildInfo()
	if ok {
		if version == "" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if commit == "" {
					commit = short(s.Value)
				}
			case "vcs.time":
				if date == "" {
					date = s.Value
				}
			}
		}
	}

	return orUnknown(version), orUnknown(commit), orUnknown(date)
}

func short(rev string) string {
	if len(rev) > 12 {
		return rev[:12]
	}
	return rev
}

func orUnknown(s string) string {
	if s == "" {
		return unknown
	}
	return s
}
