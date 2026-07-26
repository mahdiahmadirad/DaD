package cli

import (
	"runtime/debug"

	versioninfo "github.com/mahdiahmadirad/DaD/internal/version"
)

const developmentVersion = "dev"

// releaseVersion is set for release artifacts with:
//
//	-X github.com/mahdiahmadirad/DaD/internal/cli.releaseVersion=VERSION
//
// It intentionally has no release-looking default. Versioned module builds
// fall back to Go build information, while local builds report "dev".
var releaseVersion string

func CurrentVersion() string {
	moduleVersion := ""
	moduleSum := ""
	if info, ok := debug.ReadBuildInfo(); ok {
		moduleVersion = info.Main.Version
		moduleSum = info.Main.Sum
	}
	return resolveVersion(releaseVersion, moduleVersion, moduleSum)
}

func resolveVersion(injectedVersion, moduleVersion, moduleSum string) string {
	if injectedVersion != "" {
		if version, ok := versioninfo.Normalize(injectedVersion); ok {
			return version
		}
		return developmentVersion
	}

	if moduleSum != "" {
		if version, ok := versioninfo.Normalize(moduleVersion); ok {
			return version
		}
	}
	return developmentVersion
}
