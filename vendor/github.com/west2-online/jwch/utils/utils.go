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

package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func SaveData(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Output struct as json format
func PrintStruct(s interface{}) string {
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("%+v", s)
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", s)
	}

	return out.String()
}

// GetChineseCharacter returns the Chinese characters and number in the string
func GetChineseCharacter(s string) string {
	var result string
	for _, v := range s {
		if (v >= 0x4e00 && v <= 0x9fa5) || (v >= '0' && v <= '9') {
			result += string(v)
		}
	}
	return result

	// return regexp.MustCompile("[^\u4e00-\u9fa5]").ReplaceAllString(s, "")
	// After performance testing, there is no difference between the two methods
}

func RemoveDuplicate(data interface{}) interface{} {
	inArr := reflect.ValueOf(data)
	if inArr.Kind() != reflect.Slice && inArr.Kind() != reflect.Array {
		return data // 不是数组/切片
	}

	existMap := make(map[interface{}]bool)
	outArr := reflect.MakeSlice(inArr.Type(), 0, inArr.Len())

	for i := 0; i < inArr.Len(); i++ {
		iVal := inArr.Index(i)

		if _, ok := existMap[iVal.Interface()]; !ok {
			outArr = reflect.Append(outArr, inArr.Index(i))
			existMap[iVal.Interface()] = true
		}
	}

	return outArr.Interface()
}

func Base64EncodeHTTPImage(data []byte) string {
	return "data:" + http.DetectContentType(data) + "base64," + base64.StdEncoding.EncodeToString(data)
}

// Md5Hash generates MD5 hash and returns either 16 or 32 bit based on the bit parameter.
func Md5Hash(text string, bit int) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	fullHash := hex.EncodeToString(hasher.Sum(nil)) // 32-bit (full) hash

	if bit == 16 {
		return fullHash[8:24] // return 16-bit hash (substring from index 8 to 24)
	}
	return fullHash // return full 32-bit hash
}

func StructJSONEncodeBase64(data interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(PrintStruct(data)))
}

func JSONUnmarshalFromFile(filePath string, v any) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func InnerTextWithBr(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type == html.ElementNode && n.Data == "br" {
		return "\n"
	}

	var buf bytes.Buffer

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(InnerTextWithBr(c))
	}

	return buf.String()
}

func SafeAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return n
}

func ConvertGB2312ToUTF8(input []byte) (string, error) {
	// 使用 transform.NewReader 进行编码转换
	reader := transform.NewReader(bytes.NewReader(input), simplifiedchinese.GB18030.NewDecoder())

	// 将转换后的结果读取为 UTF-8 字符串
	utf8Bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(utf8Bytes), nil
}
