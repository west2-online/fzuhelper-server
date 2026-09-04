/*
Copyright 2024 The west2-online Authors.

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

package pack

import "testing"

func TestNormalizeCourseLocationPublicLanguageClassrooms(t *testing.T) {
	expectations := map[string]string{
		"旗山公语1":  "公语1（东1-508）",
		"旗山公语2":  "公语2（东1-507）",
		"旗山公语3":  "公语3（东1-506）",
		"旗山公语4":  "公语4（东1-505）",
		"旗山公语5":  "公语5（东1-504）",
		"旗山公语6":  "公语6（东1-503）",
		"旗山公语7":  "公语7（东1-502）",
		"旗山公语8":  "公语8（东1-501）",
		"旗山公语9":  "公语9（东1-408）",
		"旗山公语10": "公语10（东1-407）",
		"旗山公语11": "公语11（东1-406）",
		"旗山公语12": "公语12（东1-405）",
		"旗山公语13": "公语13（东1-404）",
		"旗山公语14": "公语14（东1-403）",
		"旗山公语15": "公语15（东1-402）",
		"旗山公语16": "公语16（东1-401）",
	}

	for location, want := range expectations {
		t.Run(location, func(t *testing.T) {
			if got := normalizeCourseLocation(location, false); got != want {
				t.Fatalf("normalizeCourseLocation(%q, false) = %q, want %q", location, got, want)
			}
		})
	}
}

func TestNormalizeCourseLocationKeepsGraduateAndSpecialLocations(t *testing.T) {
	tests := map[string]struct {
		location   string
		isGraduate bool
		want       string
	}{
		"graduate public language classroom": {
			location: "旗山公语1",
			want:     "公语1（东1-508）",
		},
		"physics laboratory center": {
			location: "旗山物理实验教学中心",
			want:     "旗山物理实验教学中心",
		},
		"copper plate teaching building": {
			location: "铜盘教学楼",
			want:     "铜盘教学楼",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := normalizeCourseLocation(tc.location, tc.isGraduate); got != tc.want {
				t.Fatalf("normalizeCourseLocation(%q, %t) = %q, want %q", tc.location, tc.isGraduate, got, tc.want)
			}
		})
	}
}
