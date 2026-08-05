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

package utils

import "testing"

func TestRemoveUndergraduatePrefix(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "9位学号-2022年入学",
			// 后9位 222200311，第3-4位 22 < 26
			id:   "20268514814222200311",
			want: "222200311",
		},
		{
			name: "9位学号-2025年入学",
			// 后9位 102512345，第3-4位 25 < 26
			id:   "20268514814102512345",
			want: "102512345",
		},
		{
			name: "10位学号-2026年入学",
			// 后9位 026012345，第3-4位 60 >= 26，取后10位
			id:   "202685148141026012345",
			want: "1026012345",
		},
		{
			name: "identifier即学号-9位",
			id:   "222200311",
			want: "222200311",
		},
		{
			name: "过短",
			id:   "12345678",
			want: "",
		},
		{
			name: "空",
			id:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveUndergraduatePrefix(tt.id); got != tt.want {
				t.Errorf("RemoveUndergraduatePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
