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

	"github.com/west2-online/fzuhelper-server/pkg/oss"
)

func TestWithOssSetProviders(t *testing.T) {
	for _, provider := range []string{oss.UpYunProvider, oss.COSProvider} {
		mockey.PatchConvey(provider, t, func() {
			upyunConfig := &oss.UpYunConfig{}
			cosConfig := &oss.CosConfig{}
			upyunCalls, cosCalls := 0, 0
			mockey.Mock(oss.NewUpYunConfig).To(func() (*oss.UpYunConfig, error) {
				upyunCalls++
				return upyunConfig, nil
			}).Build()
			mockey.Mock(oss.NewCosConfig).To(func() *oss.CosConfig {
				cosCalls++
				return cosConfig
			}).Build()

			clientSet := &ClientSet{}
			WithOssSet(provider)(clientSet)
			require.Equal(t, provider, clientSet.OssSet.Provider)
			if provider == oss.UpYunProvider {
				require.Same(t, upyunConfig, clientSet.OssSet.Upyun)
				require.Nil(t, clientSet.OssSet.Cos)
				require.Equal(t, 1, upyunCalls)
				require.Zero(t, cosCalls)
			} else {
				require.Same(t, cosConfig, clientSet.OssSet.Cos)
				require.Nil(t, clientSet.OssSet.Upyun)
				require.Equal(t, 1, cosCalls)
				require.Zero(t, upyunCalls)
			}
		})
	}
}
