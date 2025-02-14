/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package platforms

import (
	"sort"
	"testing"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Test the platform compatibility of the different
// OS Versions considering two ltsc container image
// versions (ltsc2019, ltsc2022)
func Test_PlatformCompat(t *testing.T) {
	for testName, tc := range map[string]struct {
		hostOs    uint16
		ctrOs     uint16
		shouldRun bool
	}{
		"RS5Host_ltsc2019": {
			hostOs:    rs5,
			ctrOs:     rs5,
			shouldRun: true,
		},
		"RS5Host_ltsc2022": {
			hostOs:    rs5,
			ctrOs:     v21H2Server,
			shouldRun: false,
		},
		"WS2022Host_ltsc2019": {
			hostOs:    v21H2Server,
			ctrOs:     rs5,
			shouldRun: false,
		},
		"WS2022Host_ltsc2022": {
			hostOs:    v21H2Server,
			ctrOs:     v21H2Server,
			shouldRun: true,
		},
		"Wind11Host_ltsc2019": {
			hostOs:    v22H2Win11,
			ctrOs:     rs5,
			shouldRun: false,
		},
		"Wind11Host_ltsc2022": {
			hostOs:    v22H2Win11,
			ctrOs:     v21H2Server,
			shouldRun: true,
		},
	} {
		// Check if ltsc2019/ltsc2022 guest images are compatible on
		// the given host OS versions
		//
		hostOSVersion := windowsOSVersion{
			MajorVersion: 10,
			MinorVersion: 0,
			Build:        tc.hostOs,
		}
		ctrOSVersion := windowsOSVersion{
			MajorVersion: 10,
			MinorVersion: 0,
			Build:        tc.ctrOs,
		}
		if checkWindowsHostAndContainerCompat(hostOSVersion, ctrOSVersion) != tc.shouldRun {
			var expectedResultStr string
			if !tc.shouldRun {
				expectedResultStr = " NOT"
			}
			t.Fatalf("Failed %v: host %v should%s be able to run guest %v", testName, tc.hostOs, expectedResultStr, tc.ctrOs)
		}
	}
}

func Test_PlatformOrder(t *testing.T) {
	//comparer := Only(specs.Platform{
	//	Architecture: "amd64",
	//	OS:           "windows",
	//	OSVersion:    "10.0.26100.2894",
	//	OSFeatures:   nil,
	//	Variant:      "",
	//})
	comparer := &windowsMatchComparer{Matcher: NewMatcher(specs.Platform{
		Architecture: "amd64",
		OS:           "windows",
		OSVersion:    "10.0.26100.2894",
		OSFeatures:   nil,
		Variant:      "",
	})}

	imgs := []Platform{
		{
			Architecture: "amd64",
			OS:           "linux",
			OSVersion:    "",
			OSFeatures:   nil,
			Variant:      "",
		},
		{
			Architecture: "amd64",
			OS:           "windows",
			OSVersion:    "10.0.20348.3091",
			OSFeatures:   nil,
			Variant:      "",
		},
		{
			Architecture: "amd64",
			OS:           "windows",
			OSVersion:    "10.0.26100.2894",
			OSFeatures:   nil,
			Variant:      "",
		},
	}

	sort.SliceStable(imgs, func(i, j int) bool {
		return comparer.Less(imgs[i], imgs[j])
	})

	want := "10.0.26100.2894"
	if imgs[0].OSVersion != want {
		t.Errorf("OSVersion mismatch, want %q, got %q", want, imgs[0].OSVersion)
	}
}
