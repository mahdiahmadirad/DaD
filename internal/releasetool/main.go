package main

import (
	"fmt"
	"os"
	"strings"

	releasecheck "github.com/mahdiahmadirad/DaD/internal/release"
	versioninfo "github.com/mahdiahmadirad/DaD/internal/version"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(arguments []string) error {
	if len(arguments) == 0 {
		return fmt.Errorf(
			"usage: releasetool validate TAG | notes TAG OUTPUT | verify VERSION DIST",
		)
	}
	switch arguments[0] {
	case "validate":
		if len(arguments) != 2 {
			return fmt.Errorf("usage: releasetool validate TAG")
		}
		version, valid := versioninfo.Normalize(arguments[1])
		if !valid || arguments[1] != "v"+version {
			return fmt.Errorf("%q is not a canonical release tag", arguments[1])
		}
		content, err := os.ReadFile("CHANGELOG.md")
		if err != nil {
			return err
		}
		_, err = releasecheck.ChangelogNotes(string(content), version)
		return err
	case "notes":
		if len(arguments) != 3 {
			return fmt.Errorf("usage: releasetool notes TAG OUTPUT")
		}
		content, err := os.ReadFile("CHANGELOG.md")
		if err != nil {
			return err
		}
		notes, err := releasecheck.ChangelogNotes(
			string(content), strings.TrimPrefix(arguments[1], "v"),
		)
		if err != nil {
			return err
		}
		return os.WriteFile(arguments[2], []byte(notes), 0o644)
	case "verify":
		if len(arguments) != 3 {
			return fmt.Errorf("usage: releasetool verify VERSION DIST")
		}
		return releasecheck.VerifyArtifacts(arguments[2], arguments[1])
	default:
		return fmt.Errorf("unknown release tool operation %q", arguments[0])
	}
}
