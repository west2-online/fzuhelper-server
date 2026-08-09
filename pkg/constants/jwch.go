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

package constants

const (
	MaxRetries      = 5              // 最大重试次数
	StudentIDLength = 9              // 本科基础学号长度（2026年前入学为9位），也是从 identifier 提取时的最小长度
	InitialDelay    = 1 * ONE_SECOND // 初始等待时间
	YjsTermLen      = 6              // 研究生学期长度
	JwchTermLen     = 6

	// 2026年（含）起入学的新生学号由9位变为10位。
	// 学号第3-4位为入学年份（如 22/25/26），年份 >= StudentIDYearThreshold 时为10位新学号。
	StudentIDLengthNew     = 10
	StudentIDYearThreshold = 26
)
