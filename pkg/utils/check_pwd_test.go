package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/west2-online/fzuhelper-server/config"
)

func TestCheckPwd(t *testing.T) {
	tests := []struct {
		name   string
		pwd    string
		secret string
		want   bool
	}{
		{
			name:   "Success",
			pwd:    "114514",
			secret: "114514",
			want:   true,
		},
		{
			name:   "Invalid",
			pwd:    "114514",
			secret: "1919810",
			want:   false,
		},
	}
	_ = config.InitForTest("api")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Admin.Secret = tt.secret
			result := CheckPwd(tt.pwd)
			assert.Equal(t, result, tt.want)
		})
	}
}
