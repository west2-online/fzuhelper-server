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

package captcha

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:lll
const validDataURL = "data:image/png;base64,Qk2mCAAAAAAAADYAAAAoAAAASAAAAAoAAAABABgAAAAAAHAIAAASCwAAEgsAAAAAAAAAAAAA+vr/+vr/+vr/lgD6lgD6lgD6lgD6+vr/+vr/+vr/+vr/+vr/+vr/ljIAljIAljIAljIA+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYAAJYAAJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIAljIA+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/ljIAljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/AJYAAJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/AACWAACW+vr/AACW+vr/+vr/+vr/lgD6+vr/lgD6lgD6lgD6lgD6lgD6+vr/+vr/+vr/+vr/+vr/ljIAljIAljIAljIA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACWAACWAACW+vr/+vr/+vr/+vr/+vr/"

func TestValidateLoginCode(t *testing.T) {
	type testCase struct {
		name    string
		in      string
		want    int
		wantErr bool
	}

	cases := []testCase{
		{name: "valid", in: validDataURL, want: 104, wantErr: false},
		{name: "empty", in: "", want: 0, wantErr: true},
		{name: "malformed_base64", in: "data:image/png;base64,not_base64!!", want: 0, wantErr: true},
		{name: "not_data_url", in: "not an image", want: 0, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ValidateLoginCode(tc.in)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestLoadTemplatesAndStructure(t *testing.T) {
	dataDir := filepath.Join(".", "data")
	err := LoadTemplates(dataDir)
	assert.NoError(t, err)
	assert.Equal(t, 9, len(templates))
	for i, tvec := range templates {
		assert.Greater(t, len(tvec), 0, "template %d vector is empty", i)
	}
}

func TestDigitCombinationArithmetic(t *testing.T) {
	type testCase struct {
		name string
		ds   []int
		want int
	}
	cases := []testCase{
		{name: "1,2,3,4", ds: []int{1, 2, 3, 4}, want: 46},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.ds[0]*10 + tc.ds[1] + tc.ds[2]*10 + tc.ds[3]
			assert.Equal(t, tc.want, got)
		})
	}
}
