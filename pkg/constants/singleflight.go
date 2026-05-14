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

// 要拼的动态key我放在 pkg/singleflight 里了，这里只放一些固定的 key
const (
	// 学期列表是全局校历数据，不依赖具体用户，可以使用固定 key 合并请求。
	SingleflightTermListKey = "term_list"

	// 下载地址按发布渠道拆 key，避免 release/beta 并发时复用到另一个渠道的地址。
	SingleflightDownloadReleaseKey = "download_release"
	SingleflightDownloadBetaKey    = "download_beta"

	// 版本信息按发布渠道拆 key，避免 release/beta 版本元数据互相复用。
	SingleflightReleaseVersionKey = "release_version"
	SingleflightBetaVersionKey    = "beta_version"

	SingleflightCloudKey = "cloud"

	// Android 接口一次返回 release 和 beta 两份数据，使用独立 key 避免和单独版本查询混用。
	SingleflightAndroidVersionKey = "android_version"
)
