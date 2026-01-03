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
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
)

var (
	templates [][]float64
	sites     = []int{2, 12, 32, 42}
	width     = 6
)

// init 同步加载模板。
// 如果加载失败，会在启动阶段返回错误（panic）。
func init() {
	execPath, _ := os.Getwd()
	dataDir := filepath.Join(execPath, "pkg", "captcha", "data")
	if err := LoadTemplates(dataDir); err != nil {
		panic(fmt.Errorf("pkg/captcha: load templates failed: %v", err))
	}
}

// ValidateLoginCode 识别前端提交的图片 base64，返回验证码的整数形式。
// 说明：
// - 支持 data URL（例如 "data:image/png;base64,..."），会剥离前缀并解码。
// - 使用 image.Decode 解码输入图片，然后转换为灰度并按固定位置切片匹配模板。
func ValidateLoginCode(imageString string) (int, error) {
	if imageString == "" {
		return 0, fmt.Errorf("empty image string")
	}
	if strings.Contains(imageString, ",") {
		parts := strings.SplitN(imageString, ",", 2)
		imageString = parts[1]
	}
	decoded, err := base64.StdEncoding.DecodeString(imageString)
	if err != nil {
		return 0, err
	}

	fullImg, _, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		return 0, fmt.Errorf("decode image failed: %v", err)
	}

	gray := imageToGray(fullImg)
	digits := preprocessAndRecognize(gray)
	if len(digits) != 4 {
		return 0, fmt.Errorf("recognize returned invalid length %d", len(digits))
	}
	res := digits[0]*10 + digits[1] + digits[2]*10 + digits[3]
	return res, nil
}

// preprocessAndRecognize 从整张灰度图中抽取每一位的像素向量并进行匹配。
// 关键点：
// - 该实现假定验证码四位字符在固定水平位置（由 `sites` 指定），每位宽度为 `width`。
// - 对像素采用简单二值化（阈值 250），目的是突出字符与背景，降低噪声影响。
// - 为性能考虑，按 (h*width) 预分配切片并按索引写入，避免多次扩容带来的分配开销。
func preprocessAndRecognize(gray *image.Gray) []int {
	h := gray.Bounds().Dy()
	w := gray.Bounds().Dx()
	out := make([]int, 0, len(sites))
	for _, s := range sites {
		// 如果超出边界，返回 0 作为兜底值
		if s+width > w {
			out = append(out, 0)
			continue
		}
		// 预分配像素向量（按行主序），按索引直接写入，避免 append 的多次分配
		v := make([]float64, h*width)
		idx := 0
		for y := 0; y < h; y++ {
			for x := s; x < s+width; x++ {
				p := gray.GrayAt(x, y).Y
				// 简单阈值二值化：背景接近白（>250）视为 0，否则视为 255（字符）
				if p > 250 {
					v[idx] = 0.0
				} else {
					v[idx] = 255.0
				}
				idx++
			}
		}
		out = append(out, matchIndex(v))
	}
	return out
}

// matchIndex 在已加载模板集中匹配向量 v，返回最佳模板索引。
// 实现细节：
// - 模板在加载时已经做了 L2 归一化，因此可以直接比较点积（dot）。
// - 为了兼容不同尺寸的模板和输入向量，点积只计算它们的最小长度部分。
func matchIndex(v []float64) int {
	if len(templates) == 0 {
		return 0
	}
	best := 0
	var bestDot float64
	for i, t := range templates {
		min := len(v)
		if len(t) < min {
			min = len(t)
		}
		var dot float64
		for j := 0; j < min; j++ {
			dot += v[j] * t[j]
		}
		if i == 0 || dot > bestDot {
			bestDot = dot
			best = i
		}
	}
	return best
}

// LoadTemplates 从 dataDir 中按文件名 num_0.bmp..num_8.bmp 加载模板并预计算每个模板的 L2 范数。
// 设计说明：
// - 模板使用灰度值直接作为向量分量（按行主序）。
// - 通过预计算范数避免在匹配阶段对每个模板重复计算 sqrt，提升性能。
// - 如果模板图像的范数为 0（理论上不应该出现），将其范数设为 1.0 以避免除零错误。
func LoadTemplates(dataDir string) error {
	tpls := make([][]float64, 0, 9)
	for i := 0; i < 9; i++ {
		p := filepath.Join(dataDir, fmt.Sprintf("num_%d.bmp", i))
		f, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("open template %s failed: %v", p, err)
		}
		img, err := bmp.Decode(f)
		_ = f.Close()
		if err != nil {
			return fmt.Errorf("decode bmp %s failed: %v", p, err)
		}
		b := img.Bounds()
		h := b.Dy()
		w := b.Dx()
		// 预分配固定大小的向量并按索引写入以减少内存分配与复制
		vec := make([]float64, h*w)
		var sumSquares float64
		idx := 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
				val := float64(c.Y)
				vec[idx] = val
				sumSquares += val * val
				idx++
			}
		}
		norm := math.Sqrt(sumSquares)
		if norm == 0 {
			norm = 1.0
		}
		// 归一化模板向量，便于后续直接比较点积
		for k := 0; k < len(vec); k++ {
			vec[k] = vec[k] / norm
		}
		tpls = append(tpls, vec)
	}
	templates = tpls
	return nil
}

// imageToGray 将任意图像转换为灰度图像。
// 说明：直接使用 color.GrayModel.Convert 以便兼容不同输入色彩模型。
func imageToGray(img image.Image) *image.Gray {
	b := img.Bounds()
	g := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			g.SetGray(x, y, c)
		}
	}
	return g
}
