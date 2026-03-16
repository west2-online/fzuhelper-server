// Copyright 2023 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"sync"
	"sync/atomic"

	"github.com/bytedance/sonic"

	"github.com/cloudwego/configmanager/iface"
)

var counter int64 = 0

// sonicAPI is used for **decoding** json with UseNumber option.
var sonicAPI = sonic.Config{UseNumber: true}.Froze()

// IncreaseCounter increases the counter by 1 and returns the new value.
// It's useful when you need a process-wide unique ID.
func IncreaseCounter() int64 {
	return atomic.AddInt64(&counter, 1)
}

// JSONSerializer is a configuration serializer that encodes and decodes data in JSON format
var JSONSerializer iface.ConfigSerializer = &configSerializerJSON{}

// CompareConfig compares the config changes between the current and latest version of the configuration.
// It takes in two maps that contain the configurations as key-value pairs and returns a slice of
// ConfigChange which represent the differences between the two versions.
func CompareConfig(current, latest map[string]iface.ConfigValue) []iface.ConfigChange {
	differences := make([]iface.ConfigChange, 0)
	for k, cv := range current {
		if lv, ok := latest[k]; ok {
			// update
			if !cv.EqualsTo(lv) {
				differences = append(differences, iface.ConfigChange{
					Type:      iface.ConfigChangeTypeUpdate,
					ConfigKey: k,
					OldValue:  cv,
					NewValue:  lv,
				})
			}
		} else {
			// delete
			differences = append(differences, iface.ConfigChange{
				Type:      iface.ConfigChangeTypeDelete,
				ConfigKey: k,
				OldValue:  cv,
				NewValue:  nil,
			})
		}
	}
	for k, lv := range latest {
		if _, ok := current[k]; !ok {
			// add
			differences = append(differences, iface.ConfigChange{
				Type:      iface.ConfigChangeTypeAdd,
				ConfigKey: k,
				OldValue:  nil,
				NewValue:  lv,
			})
		}
	}
	return differences
}

// ConfigSyncMapToMap is used to convert a sync.Map to a map[string]iface.ConfigValue.
func ConfigSyncMapToMap(m *sync.Map) map[string]iface.ConfigValue {
	result := make(map[string]iface.ConfigValue) // always return a non-nil map
	if m != nil {
		m.Range(func(key, value interface{}) bool {
			result[key.(string)] = value.(iface.ConfigValue)
			return true
		})
	}
	return result
}

// ConfigMapToSyncMap returns a new sync.Map object converted from a map[string]iface.ConfigValue.
func ConfigMapToSyncMap(m map[string]iface.ConfigValue) *sync.Map {
	result := &sync.Map{} // always return a non-nil sync.Map
	for key, value := range m {
		result.Store(key, value)
	}
	return result
}

type configSerializerJSON struct{}

// Encode serializes a map of configurations into JSON format.
// It returns a byte slice containing the encoded data and an error if it fails to encode.
func (s *configSerializerJSON) Encode(config map[string]iface.ConfigValue) ([]byte, error) {
	return sonic.ConfigStd.MarshalIndent(config, "", "\t")
}

// JsonDecodeWithNumber decodes a JSON byte slice with the option to parse numbers as json.Number type.
// This function is useful when precision is important and dealing with large values.
func JsonDecodeWithNumber(b []byte, v interface{}) error {
	return sonicAPI.Unmarshal(b, v)
}

// MustJsonMarshal marshals the given interface into a JSON string.
// If the operation fails, it panics.
func MustJsonMarshal(v interface{}) []byte {
	data, err := sonic.ConfigStd.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// MustJsonUnmarshal unmarshal the given json bytes into given interface; panics on failure.
func MustJsonUnmarshal(b []byte, v interface{}) {
	if err := sonicAPI.Unmarshal(b, v); err != nil {
		panic(err)
	}
}

// JsonInitializer receives an allocator func which produces an empty Item object, and returns a function which
// fills a new Item object (by the allocator) with given json bytes.
func JsonInitializer(allocator func() iface.ConfigValueItem) iface.ItemInitializer {
	return func(js []byte) (iface.ConfigValueItem, error) {
		i := allocator()
		if err := sonicAPI.Unmarshal(js, i); err == nil {
			return i, nil
		} else {
			return nil, err
		}
	}
}
