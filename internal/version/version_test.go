package version

import "testing"

func TestNormalize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
		valid bool
	}{
		{input: "0.1.0", want: "0.1.0", valid: true},
		{input: "v1.2.3", want: "1.2.3", valid: true},
		{input: "v1.2.3-rc.1", want: "1.2.3-rc.1", valid: true},
		{input: "1.2.3+build.7", want: "1.2.3+build.7", valid: true},
		{
			input: "v0.0.0-20260727000000-deadbeef1234",
			want:  "0.0.0-20260727000000-deadbeef1234",
			valid: true,
		},
		{
			input: "123456789012345678901234567890.2.3",
			want:  "123456789012345678901234567890.2.3",
			valid: true,
		},
		{input: "", valid: false},
		{input: "(devel)", valid: false},
		{input: "dev", valid: false},
		{input: "1.2", valid: false},
		{input: "01.2.3", valid: false},
		{input: "1.2.3-01", valid: false},
		{input: "1.2.3-", valid: false},
		{input: "1.2.3+bad!", valid: false},
		{input: " 1.2.3", valid: false},
	}
	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			got, valid := Normalize(test.input)
			if got != test.want || valid != test.valid {
				t.Fatalf(
					"Normalize(%q) = (%q, %t), want (%q, %t)",
					test.input, got, valid, test.want, test.valid,
				)
			}
		})
	}
}
