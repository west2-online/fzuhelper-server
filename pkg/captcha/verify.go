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
	"runtime"
	"strings"

	"golang.org/x/image/bmp"
)

var (
	templates [][]float64
	sites     = []int{2, 12, 32, 42}
	width     = 6
)

const (
	dataURLParts  = 2
	captchaDigits = 4
	binThreshold  = 250
	templateCount = 9
)

// init 同步加载模板。
// 如果加载失败，会在启动阶段返回错误（panic）。
func init() {
	// 尝试通过源文件位置推断 data 目录，避免依赖运行时的工作目录。
	// 这样在运行 `go test` 或从不同工作目录执行时都能正确找到模板数据。
	_, filename, _, ok := runtime.Caller(0)
	var dataDir string
	if ok {
		dataDir = filepath.Join(filepath.Dir(filename), "data")
	} else {
		// 兜底回退到当前工作目录下的相对路径（历史行为），以提升兼容性
		execPath, _ := os.Getwd()
		dataDir = filepath.Join(execPath, "pkg", "captcha", "data")
	}
	if err := LoadTemplates(dataDir); err != nil {
		panic(fmt.Errorf("pkg/captcha: load templates failed: %w", err))
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
		parts := strings.SplitN(imageString, ",", dataURLParts)
		imageString = parts[1]
	}
	decoded, err := base64.StdEncoding.DecodeString(imageString)
	if err != nil {
		return 0, err
	}

	fullImg, _, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		return 0, fmt.Errorf("decode image failed: %w", err)
	}

	gray := imageToGray(fullImg)
	digits := preprocessAndRecognize(gray)
	if len(digits) != captchaDigits {
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
				// 简单阈值二值化：背景接近白（>binThreshold）视为 0，否则视为 255（字符）
				if p > binThreshold {
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
// - 模板以原始灰度向量存储；匹配时会对输入向量和模板在重叠维度上同时进行 L2 归一化并计算余弦相似度。
// - 为了兼容不同尺寸的模板和输入向量，计算只使用它们的最小长度部分。
func matchIndex(v []float64) int {
	if len(templates) == 0 {
		return 0
	}
	sims := getCosSimilarMulti(v, templates)
	best := 0
	var bestSim float64
	for i, s := range sims {
		if i == 0 || s > bestSim {
			bestSim = s
			best = i
		}
	}
	return best
}

// getCosSimilarMulti 计算向量 v 与多个模板的余弦相似度：
// - 对每个模板使用重叠长度计算点积与范数
// - 对 0 范数用 eps 或 1.0 防止除零
// - 结果映射为 0.5 + 0.5 * cos
func getCosSimilarMulti(v []float64, tpl [][]float64) []float64 {
	res := make([]float64, 0, len(tpl))
	for _, t := range tpl {
		min := len(v)
		if len(t) < min {
			min = len(t)
		}
		var dot float64
		var vSumSquares float64
		var tSumSquares float64
		for j := 0; j < min; j++ {
			dot += v[j] * t[j]
			vSumSquares += v[j] * v[j]
			tSumSquares += t[j] * t[j]
		}
		denom := math.Sqrt(vSumSquares * tSumSquares)
		if denom == 0 {
			denom = 1.0
		}
		cos := dot / denom
		if math.IsInf(cos, -1) || math.IsNaN(cos) {
			cos = 0
		}
		sim := 0.5 + 0.5*cos
		res = append(res, sim)
	}
	return res
}

// LoadTemplates 从 dataDir 中按文件名 num_0.bmp..num_8.bmp 加载模板为原始灰度向量（按行主序）。
// 设计说明：
// - 模板以原始像素值存储；匹配阶段会按重叠长度计算余弦相似度并做归一化。
// - 如果模板图像的范数为 0（理论上不应该出现），匹配计算中会采取措施避免除零。
func LoadTemplates(dataDir string) error {
	tpls := make([][]float64, 0, templateCount)
	for i := 0; i < templateCount; i++ {
		p := filepath.Join(dataDir, fmt.Sprintf("num_%d.bmp", i))
		f, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("open template %s failed: %w", p, err)
		}
		img, err := bmp.Decode(f)
		_ = f.Close()
		if err != nil {
			return fmt.Errorf("decode bmp %s failed: %w", p, err)
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
				col := color.GrayModel.Convert(img.At(x, y))
				if cg, ok := col.(color.Gray); ok {
					val := float64(cg.Y)
					vec[idx] = val
					sumSquares += val * val
				} else {
					vec[idx] = 0
				}
				idx++
			}
		}
		_ = math.Sqrt(sumSquares)
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
			col := color.GrayModel.Convert(img.At(x, y))
			if cg, ok := col.(color.Gray); ok {
				g.SetGray(x, y, cg)
			} else {
				// fallback to computed grayscale value
				r, gcol, bcol, _ := img.At(x, y).RGBA()
				yv := uint8((r*299 + gcol*587 + bcol*114) / 1000 >> 8)
				g.SetGray(x, y, color.Gray{Y: yv})
			}
		}
	}
	return g
}
