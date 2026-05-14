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

package singleflight

import "testing"

func TestDynamicKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "scores", got: ScoresKey("0001", true), want: "scores:0001:true"},
		{name: "exam rooms", got: ExamRoomsKey("0001", "202401", false), want: "exam_rooms:0001:202401:false"},
		{name: "course list", got: CourseListKey("0001", "202401", true, false), want: "courses:0001:202401:true:false"},
		{name: "course terms", got: CourseTermsKey("0001", false), want: "terms:0001:false"},
		{name: "term", got: TermKey("202401"), want: "term:202401"},
		{name: "notice", got: NoticeKey(2), want: "notice:2"},
		{name: "paper dir", got: PaperDirKey("/foo"), want: "dir:/foo"},
		{name: "user info", got: UserInfoKey("0001", true), want: "user_info:0001:true"},
		{name: "friend list", got: FriendListKey("0001"), want: "friend_list:0001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.got != tt.want {
				t.Fatalf("key = %q, want %q", tt.got, tt.want)
			}
		})
	}
}
