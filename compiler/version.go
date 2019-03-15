package compiler

import "fmt"

const (
	// VersionMajor is the Major version component of the current release
	VersionMajor = 0
	// VersionMinor is the Minor version component of the current release
	VersionMinor = 1
	// VersionPatch is the Patch version component of the current release
	VersionPatch = 1
)

// Git SHA1 commit hash of the release (set via linker flags)
var GitCommit = ""

// Version holds the textual version string.
var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()

// VersionWithCommit holds the textual version and the first 8 character of git commit.
func VersionWithCommit(gitCommit string) string {
	version := Version
	if len(gitCommit) >= 8 {
		version += "+" + gitCommit[:8]
	}
	return version
}
