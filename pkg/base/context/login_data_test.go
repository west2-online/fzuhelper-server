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

package context

import "testing"

func TestExtractIDFromIdentifier(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "研究生-9位学号",
			// 00000 + 102301001
			id:   "00000102301001",
			want: "102301001",
		},
		{
			name: "研究生-10位学号",
			// 00000 + 1026012345
			id:   "000001026012345",
			want: "1026012345",
		},
		{
			name: "本科生-9位学号",
			id:   "20268514814222200311",
			want: "222200311",
		},
		{
			name: "本科生-10位学号",
			id:   "202685148140102612345",
			want: "0102612345",
		},
		{
			name: "过短",
			id:   "123",
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
			if got := ExtractIDFromIdentifier(tt.id); got != tt.want {
				t.Errorf("ExtractIDFromIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
