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
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"math"
	"strings"
	"sync/atomic"

	"golang.org/x/image/bmp"
)

var (
	templates atomic.Value           // stores [][]float64 模板集，每个模板为原始灰度向量
	sites     = []int{2, 12, 32, 42} // 每位验证码字符在图片中的水平起始位置
	width     = 6                    // 每位验证码字符的宽度
)

const (
	// maxImageSize 图片 base64 字符串的最大长度(1MB)
	// 正常验证码图片 base64 编码后通常只有几 KB，1MB 是一个安全的上限
	maxImageSize = 1 << 20 // 1MB = 1048576 bytes
)

// Init 同步加载模板(从内置 base64 数据)。
// 必须在使用验证码功能前显式调用,如果加载失败会返回错误。
func Init() error {
	if err := loadTemplates(); err != nil {
		return fmt.Errorf("captcha: load templates failed: %w", err)
	}
	return nil
}

// ValidateLoginCode 识别前端提交的图片 base64，返回验证码的整数形式。
// 说明：
// - 支持 data URL（例如 "data:image/png;base64,..."），会剥离前缀并解码。
// - 使用 image.Decode 解码输入图片，然后转换为灰度并按固定位置切片匹配模板。
func ValidateLoginCode(imageString string) (int, error) {
	if imageString == "" {
		return 0, fmt.Errorf("empty image string")
	}
	// 防止超大图片攻击
	if len(imageString) > maxImageSize {
		return 0, fmt.Errorf("image string too large: %d bytes (max: %d)", len(imageString), maxImageSize)
	}
	// 剥离 data URL 前缀
	_, imageString, _ = strings.Cut(imageString, ",")
	decoded, err := base64.StdEncoding.DecodeString(imageString)
	if err != nil {
		return 0, err
	}
	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		return 0, fmt.Errorf("decode image failed: %w", err)
	}
	// 转为灰度图并识别
	gray := convertImageToGray(img)
	digits := recognizeImage(gray)
	// 计算结果必须是4个数字
	if len(digits) != 4 { //nolint:mnd
		return 0, fmt.Errorf("recognize returned invalid length %d", len(digits))
	}
	res := digits[0]*10 + digits[1] + digits[2]*10 + digits[3] // 验证码结果为两位数加两位数
	return res, nil
}

// recognizeImage 从整张灰度图中抽取每一位的像素向量并进行匹配。
// 关键点：
// - 该实现假定验证码四位字符在固定水平位置（由 `sites` 指定），每位宽度为 `width`。
// - 先对整张图像做二值化（阈值 250，反转），然后提取每位字符区域并匹配。
func recognizeImage(gray *image.Gray) []int {
	// 获取模板
	tplsIface := templates.Load()
	if tplsIface == nil {
		return make([]int, len(sites))
	}
	tpls, ok := tplsIface.([][]float64)
	if !ok || len(tpls) == 0 {
		return make([]int, len(sites))
	}
	// 对整张图像进行二值化阈值处理（>250 -> 0, <=250 -> 255）
	h := gray.Bounds().Dy()
	w := gray.Bounds().Dx()
	for i := 0; i < len(gray.Pix); i++ {
		// 阈值250，大于则为背景
		if gray.Pix[i] > 250 { //nolint:mnd
			gray.Pix[i] = 0
		} else {
			gray.Pix[i] = 255
		}
	}
	// 提取每个字符位置并匹配
	out := make([]int, len(sites))
	for idx, s := range sites {
		// 如果超出边界，返回 0 作为兜底值
		if s+width > w {
			out[idx] = 0
			continue
		}
		// 提取字符区域像素（所有行，列从s到s+width）
		v := make([]float64, h*width)
		vidx := 0
		for y := range h {
			for x := s; x < s+width; x++ {
				v[vidx] = float64(gray.Pix[y*gray.Stride+x])
				vidx++
			}
		}
		// 计算与所有模板的相似度，找最大值
		sims := getCosSimilarMulti(v, tpls)
		best := 0
		for i := 1; i < len(sims); i++ {
			if sims[i] > sims[best] {
				best = i
			}
		}
		out[idx] = best
	}
	return out
}

// getCosSimilarMulti 计算向量 v 与多个模板的余弦相似度：
// - 对每个模板使用重叠长度计算点积与范数
// - 结果映射为 0.5 + 0.5 * cos
func getCosSimilarMulti(v []float64, tpl [][]float64) []float64 {
	res := make([]float64, len(tpl))
	// 预先计算 v 的范数平方
	var vNormSq float64
	for _, val := range v {
		vNormSq += val * val
	}
	// 逐个计算与模板的余弦相似度
	for i, t := range tpl {
		// 计算重叠长度
		minLen := min(len(v), len(t))
		// 计算点积和模板范数平方
		var dot, tNormSq float64
		for j := range minLen {
			dot += v[j] * t[j]
			tNormSq += t[j] * t[j]
		}
		// 计算余弦相似度
		denomSq := vNormSq * tNormSq
		if denomSq == 0 {
			res[i] = 0.0
			continue
		}
		cos := dot / math.Sqrt(denomSq)
		// 将 [-1,1] 映射到 [0,1]
		res[i] = 0.5 + 0.5*cos //nolint:mnd
	}
	return res
}

// loadTemplates 从 data.go 内置的 base64 BMP 字符串加载模板为原始灰度向量。
func loadTemplates() error {
	base64Templates := GetTemplatesData()
	tpls := make([][]float64, len(base64Templates))
	for i, tplB64 := range base64Templates {
		decoded, err := base64.StdEncoding.DecodeString(tplB64)
		if err != nil {
			return fmt.Errorf("decode template %d base64 failed: %w", i, err)
		}
		img, err := bmp.Decode(bytes.NewReader(decoded))
		if err != nil {
			return fmt.Errorf("decode template %d bmp failed: %w", i, err)
		}
		// 转为灰度图并展平像素数组
		gray := convertImageToGray(img)
		vec := make([]float64, len(gray.Pix))
		for j, p := range gray.Pix {
			vec[j] = float64(p)
		}
		tpls[i] = vec
	}
	templates.Store(tpls)
	return nil
}

// convertImageToGray 将任意图像转换为灰度图像，利用 draw.Draw 做颜色转换。
func convertImageToGray(img image.Image) *image.Gray {
	b := img.Bounds()                     // 获取图像边界
	g := image.NewGray(b)                 // 创建一个新的灰度图像
	draw.Draw(g, b, img, b.Min, draw.Src) // 将原图像绘制到灰度图像上
	return g
}
