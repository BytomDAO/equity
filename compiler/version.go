package compiler

import "fmt"

const (
	VersionMajor = 0 // Major version component of the current release
	VersionMinor = 1 // Minor version component of the current release
	VersionPatch = 1 // Patch version component of the current release
)

// Git SHA1 commit hash of the release (set via linker flags)
var GitCommit = ""

// Version holds the textual version string.
var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()

func VersionWithCommit(gitCommit string) string {
	version := Version
	if len(gitCommit) >= 8 {
		version += "+" + gitCommit[:8]
	}
	return version
}
