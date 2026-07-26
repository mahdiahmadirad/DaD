package cli

import "testing"

func TestCurrentVersionRejectsInvalidInjectedVersion(t *testing.T) {
	previous := releaseVersion
	releaseVersion = "not-a-version"
	t.Cleanup(func() {
		releaseVersion = previous
	})

	if got := CurrentVersion(); got != developmentVersion {
		t.Fatalf("CurrentVersion() = %q, want %q", got, developmentVersion)
	}
}

func TestResolveVersionSources(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		injected  string
		module    string
		moduleSum string
		want      string
	}{
		{
			name:      "release injection",
			injected:  "1.2.3",
			module:    "v9.9.9",
			moduleSum: "h1:module",
			want:      "1.2.3",
		},
		{
			name:      "versioned module",
			module:    "v1.2.3",
			moduleSum: "h1:module",
			want:      "1.2.3",
		},
		{
			name:      "invalid injection does not fall back",
			injected:  "invalid",
			module:    "v1.2.3",
			moduleSum: "h1:module",
			want:      developmentVersion,
		},
		{
			name:   "local VCS pseudo-version",
			module: "v0.0.0-20260726214505-77497fe8d33d+dirty",
			want:   developmentVersion,
		},
		{
			name:   "development module",
			module: "(devel)",
			want:   developmentVersion,
		},
		{
			name: "missing build information",
			want: developmentVersion,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := resolveVersion(
				test.injected, test.module, test.moduleSum,
			); got != test.want {
				t.Fatalf(
					"resolveVersion(%q, %q, %q) = %q, want %q",
					test.injected, test.module, test.moduleSum, got, test.want,
				)
			}
		})
	}
}
