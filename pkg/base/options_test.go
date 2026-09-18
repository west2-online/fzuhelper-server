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

package base

import (
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/require"
	tencentyun "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/cos"
	"github.com/west2-online/fzuhelper-server/pkg/oss"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestWithOssSetCOSProvider(t *testing.T) {
	mockey.PatchConvey(oss.COSProvider, t, func() {
		cosConfig := &oss.CosConfig{}
		mockey.Mock(oss.NewCosConfig).Return(cosConfig).Build()
		clientSet := &ClientSet{}
		WithOssSet(oss.COSProvider)(clientSet)
		require.Equal(t, oss.COSProvider, clientSet.OssSet.Provider)
		require.Same(t, cosConfig, clientSet.OssSet.Cos)
	})
}

func TestWithFeedbackCOSClient(t *testing.T) {
	mockey.PatchConvey("OA COS config and file names", t, func() {
		original := config.Cos
		t.Cleanup(func() { config.Cos = original })
		require.NoError(t, config.InitForTest("oa"))
		require.NotNil(t, config.Cos)
		require.Equal(t, "/feedback/", config.Cos.Path)
		mockey.Mock(cos.NewCos).Return(&tencentyun.Client{}).Build()
		sf, err := utils.NewSnowflake(0, 0)
		require.NoError(t, err)
		clientSet := &ClientSet{SFClient: sf}
		WithFeedbackCOSClient()(clientSet)
		require.NotNil(t, clientSet.FeedbackCOSClient)
		url, remotePath, err := clientSet.FeedbackCOSClient.GenerateFileName(cos.FeedbackLogCategory, cos.FeedbackLogFileExtension)
		require.NoError(t, err)
		require.Regexp(t, `^/feedback/log/[0-9]+\.json\.gz$`, remotePath)
		require.Equal(t, config.Cos.DownloadDomain+remotePath, url)
	})
}
