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

package iface

import "errors"

// ConfigChangeType is an enumeration that represents the different types of config changes that can occur.
// Possible values are ConfigChangeTypeAdd, ConfigChangeTypeUpdate, and ConfigChangeTypeDelete.
type ConfigChangeType = int

const (
	ConfigChangeTypeAdd    ConfigChangeType = 1
	ConfigChangeTypeUpdate ConfigChangeType = 2
	ConfigChangeTypeDelete ConfigChangeType = 3
)

var (
	ErrNotModified       = errors.New("not modified")
	ErrWaitTimeout       = errors.New("wait timeout")
	ErrConfigNotFound    = errors.New("not exists")
	ErrInvalidConfigKey  = errors.New("invalid config key")
	ErrItemNotFound      = errors.New("item not found")
	ErrNoItemInitializer = errors.New("item initializer not available")
)

// LogFunc is used to decouple config manager with logger implementations
type LogFunc func(format string, args ...interface{})

// ConfigValue is business related, and config manager don't care about its content.
type ConfigValue interface {
	DeepCopy() ConfigValue
	EqualsTo(ConfigValue) bool
	GetItem(itemType ItemType) (ConfigValueItem, error)
	GetItemOrDefault(itemType ItemType, defaultItem ConfigValueItem) ConfigValueItem
}

// ItemType is the type of unique identifier for ConfigValueItem
type ItemType string

// ConfigValueItem interface represents an item in a configuration value
type ConfigValueItem interface {
	DeepCopy() ConfigValueItem
	EqualsTo(ConfigValueItem) bool
}

// ItemInitializer defines a function that takes in a byte slice and returns a ConfigValueItem and an error.
type ItemInitializer func([]byte) (ConfigValueItem, error)

// ConfigProvider should implement the logic of loading config from somewhere, for example, a local file,
// or a remote config center
type ConfigProvider interface {
	// LoadConfig should return a full set of config, or an error if failed. currentConfig is given as a
	// 2nd parameter so that the Not-Modified strategy of some provider (if applicable) can be used.
	// If the config is not modified at all, (nil, ErrNotModified) as a quick response is allowed.
	LoadConfig(currentConfig map[string]ConfigValue) (map[string]ConfigValue, error)
}

// ConfigChange defines the change of a ConfigValue.
type ConfigChange struct {
	Type      ConfigChangeType
	ConfigKey string
	OldValue  ConfigValue
	NewValue  ConfigValue
}

// ConfigChangeListener defines the callback function when a ConfigValue is changed.
type ConfigChangeListener func(change ConfigChange)

// RefreshCompleteSignal receives a signal when the refresh is done.
type RefreshCompleteSignal = chan error

// ConfigManagerIface defines the interface of a Config Manager
type ConfigManagerIface interface {
	GetProvider() ConfigProvider
	RegisterConfigChangeListener(identifier string, listener ConfigChangeListener)
	DeregisterConfigChangeListener(identifier string)
	Refresh() RefreshCompleteSignal
	RefreshAndWait() error
	GetConfig(key string) (ConfigValue, error)
	GetConfigItem(key string, itemType ItemType) (ConfigValueItem, error)
	GetAllConfig() map[string]ConfigValue
	Dump(filepath string) error
}

// ConfigSerializer is an interface that defines methods for serializing and deserializing configuration data.
type ConfigSerializer interface {
	Encode(config map[string]ConfigValue) ([]byte, error)
}

// Limiter is an interface that defines methods for limiting the rate of actions.
type Limiter interface {
	Allow() bool
}
