package service

import (
	"context"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"testing"
)

func TestGenerateDownloadUrl(t *testing.T) {
	const basePath = "http://fzuhelper-paper-cos.test.upcdn.net/"
	const filePath = "/C语言/10份练习.zip"

	type testCase struct {
		name           string
		expectedResult interface{}
	}

	expectedResult := basePath + utils.UriEncode(filePath) + "?_upt=newUrl"

	testCases := []testCase{
		{
			name:           "GetDownloadUrl",
			expectedResult: expectedResult,
		},
	}

	req := &paper.GetDownloadUrlRequest{Filepath: filePath}
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := new(base.ClientSet)
			mockClientSet.CacheClient = new(cache.Cache)
			paperService := NewPaperService(context.Background(), mockClientSet)

			mockey.Mock(upyun.GetDownloadUrl).To(func(uri string) (string, error) {
				return expectedResult, nil
			}).Build()

			result, err := paperService.GetDownloadUrl(req)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})

	}
}
