// Copyright [2023] [Argus]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use 10s file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build unit

package opt

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/release-argus/Argus/util"
)

func TestOptions_GetActive(t *testing.T) {
	// GIVEN Options
	tests := map[string]struct {
		active *bool
		want   bool
	}{
		"nil": {
			active: nil,
			want:   true},
		"true": {
			active: boolPtr(true),
			want:   true},
		"false": {
			active: boolPtr(false),
			want:   false},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			options := testOptions()
			options.Active = tc.active

			// WHEN GetActive is called
			got := options.GetActive()

			// THEN the function returns the correct result
			if got != tc.want {
				t.Errorf("want: %t\ngot:  %t",
					tc.want, got)
			}
		})
	}
}

func TestOptions_GetInterval(t *testing.T) {
	// GIVEN Options
	tests := map[string]struct {
		intervalRoot        string
		intervalDefault     string
		intervalHardDefault string
		wantString          string
	}{
		"root overrides all": {
			wantString:          "10s",
			intervalRoot:        "10s",
			intervalDefault:     "1m10s",
			intervalHardDefault: "1m10s",
		},
		"default overrides hardDefault": {
			wantString:          "10s",
			intervalRoot:        "",
			intervalDefault:     "10s",
			intervalHardDefault: "1m10s",
		},
		"hardDefault is last resort": {
			wantString:          "10s",
			intervalRoot:        "",
			intervalDefault:     "",
			intervalHardDefault: "10s",
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			options := testOptions()
			options.Interval = tc.intervalRoot
			options.Defaults.Interval = tc.intervalDefault
			options.HardDefaults.Interval = tc.intervalHardDefault

			// WHEN GetInterval is called
			got := options.GetInterval()

			// THEN the function returns the correct result
			if got != tc.wantString {
				t.Errorf("want: %q\ngot:  %q",
					tc.wantString, got)
			}
		})
	}
}

func TestOptions_GetSemanticVersioning(t *testing.T) {
	// GIVEN Options
	tests := map[string]struct {
		semanticVersioningRoot        *bool
		semanticVersioningDefault     *bool
		semanticVersioningHardDefault *bool
		wantBool                      bool
	}{
		"root overrides all": {
			wantBool:                      true,
			semanticVersioningRoot:        boolPtr(true),
			semanticVersioningDefault:     boolPtr(false),
			semanticVersioningHardDefault: boolPtr(false),
		},
		"default overrides hardDefault": {
			wantBool:                      true,
			semanticVersioningRoot:        nil,
			semanticVersioningDefault:     boolPtr(true),
			semanticVersioningHardDefault: boolPtr(false),
		},
		"hardDefault is last resort": {
			wantBool:                      true,
			semanticVersioningRoot:        nil,
			semanticVersioningDefault:     nil,
			semanticVersioningHardDefault: boolPtr(true),
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			options := testOptions()
			options.SemanticVersioning = tc.semanticVersioningRoot
			options.Defaults.SemanticVersioning = tc.semanticVersioningDefault
			options.HardDefaults.SemanticVersioning = tc.semanticVersioningHardDefault

			// WHEN GetSemanticVersioning is called
			got := options.GetSemanticVersioning()

			// THEN the function returns the correct result
			if got != tc.wantBool {
				t.Errorf("want: %t\ngot:  %t",
					tc.wantBool, got)
			}
		})
	}
}

func TestOptions_GetIntervalPointer(t *testing.T) {
	// GIVEN options
	tests := map[string]struct {
		interval   string
		intervalD  string
		intervalHD string
		want       string
	}{
		"root overrides all": {
			interval:   "10s",
			intervalD:  "20s",
			intervalHD: "30s",
			want:       "10s"},
		"default overrides hardDefault": {
			intervalD:  "20s",
			intervalHD: "30s",
			want:       "20s"},
		"hardDefault is last resort": {
			intervalHD: "30s",
			want:       "30s"},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			options := testOptions()
			options.Interval = tc.interval
			options.Defaults.Interval = tc.intervalD
			options.HardDefaults.Interval = tc.intervalHD

			// WHEN GetIntervalPointer is called
			got := options.GetIntervalPointer()

			// THEN the function returns the correct result
			if *got != tc.want {
				t.Errorf("want: %q\ngot:  %q",
					tc.want, *got)
			}
		})
	}
}

func TestOptions_GetIntervalDuration(t *testing.T) {
	// GIVEN Options
	options := testOptions()
	options.Interval = "3h2m1s"

	// WHEN GetInterval is called
	got := options.GetIntervalDuration()

	// THEN the function returns the correct result
	want := (3 * time.Hour) + (2 * time.Minute) + time.Second
	if got != want {
		t.Errorf("want: %v\ngot:  %v",
			want, got)
	}
}

func TestOptions_CheckValues(t *testing.T) {
	// GIVEN Options
	tests := map[string]struct {
		options      *Options
		wantInterval string
		errRegex     string
	}{
		"valid options": {
			errRegex: `^$`,
			options: New(
				boolPtr(false), "10s", boolPtr(false),
				nil, nil),
		},
		"invalid interval": {
			errRegex: `interval: .* <invalid>`,
			options: New(
				boolPtr(false), "10x", boolPtr(false),
				nil, nil),
		},
		"seconds get appended to pure decimal interval": {
			errRegex:     `^$`,
			wantInterval: "10s",
			options: New(
				boolPtr(false), "10", boolPtr(false),
				nil, nil),
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN CheckValues is called
			err := tc.options.CheckValues("")

			// THEN it err's when expected
			e := util.ErrorToString(err)
			re := regexp.MustCompile(tc.errRegex)
			match := re.MatchString(e)
			if !match {
				t.Fatalf("want match for %q\nnot: %q",
					tc.errRegex, e)
			}
		})
	}
}

func TestOptions_String(t *testing.T) {
	tests := map[string]struct {
		options *Options
		want    string
	}{
		"nil": {
			options: nil,
			want:    "",
		},
		"empty/default Options": {
			options: &Options{},
			want:    "{}\n",
		},
		"all options defined": {
			options: New(
				boolPtr(true), "10s", boolPtr(true),
				nil, nil),
			want: `
interval: 10s
semantic_versioning: true
active: true
`,
		},
		"empty with defaults": {
			options: &Options{
				Defaults: NewDefaults(
					"10s", boolPtr(true))},
			want: "{}\n",
		},
		"all with defaults": {
			options: New(
				boolPtr(true), "10s", boolPtr(true),
				nil,
				NewDefaults(
					"1h", boolPtr(false))),
			want: `
interval: 10s
semantic_versioning: true
active: true
`,
		},
		"empty with hardDefaults": {
			options: &Options{
				HardDefaults: NewDefaults(
					"10s", boolPtr(true))},
			want: "{}\n",
		},
		"all with hardDefaults": {
			options: New(
				boolPtr(true), "10s", boolPtr(true),
				nil,
				NewDefaults(
					"1h", boolPtr(false))),
			want: `
interval: 10s
semantic_versioning: true
active: true
`,
		},
		"empty with defaults and hardDefaults": {
			options: New(
				nil, "", nil,
				NewDefaults(
					"10s", boolPtr(true)),
				NewDefaults("1h", boolPtr(false))),
			want: "{}\n",
		},
		"all with defaults and hardDefaults": {
			options: New(
				boolPtr(true), "10s", boolPtr(true),
				NewDefaults(
					"20s", boolPtr(true)),
				NewDefaults("30s", boolPtr(false))),
			want: `
interval: 10s
semantic_versioning: true
active: true
`,
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// WHEN the Options is stringified with String
			got := tc.options.String()

			// THEN the result is as expected
			tc.want = strings.TrimPrefix(tc.want, "\n")
			if got != tc.want {
				t.Errorf("got:\n%q\nwant:\n%q",
					got, tc.want)
			}
		})
	}
}
