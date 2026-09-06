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

package common

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/internal/common/service"
	kitexCommon "github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/jwch"
)

func TestCommonServiceImplGetTermTermFormat(t *testing.T) {
	type TestCase struct {
		Name       string
		reqTerm    string
		expectTerm string
	}

	testCases := []TestCase{
		{
			Name:       "YjsyTermConvertedToJwchTerm",
			reqTerm:    "2026-2027-1",
			expectTerm: "202601",
		},
		{
			Name:       "YjsySecondSemesterConvertedToJwchTerm",
			reqTerm:    "2024-2025-2",
			expectTerm: "202402",
		},
		{
			Name:       "JwchTermUnchanged",
			reqTerm:    "202601",
			expectTerm: "202601",
		},
	}

	events := &jwch.CalTermEvents{
		TermId:     "202601",
		Term:       "202601",
		SchoolYear: "2026",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			var capturedTerm string
			mockey.Mock((*service.CommonService).GetTerm).To(
				func(_ *service.CommonService, req *kitexCommon.TermRequest) (bool, *jwch.CalTermEvents, error) {
					capturedTerm = req.Term
					return true, events, nil
				},
			).Build()

			s := NewCommonService(&base.ClientSet{}, nil)
			resp, err := s.GetTerm(context.Background(), &kitexCommon.TermRequest{Term: tc.reqTerm})

			assert.Nil(t, err)
			assert.Equal(t, tc.expectTerm, capturedTerm)
			assert.NotNil(t, resp)
			assert.Equal(t, "202601", resp.TermInfo.GetTerm())
			assert.Equal(t, int64(errno.SuccessCode), resp.Base.GetCode())
		})
	}
}
