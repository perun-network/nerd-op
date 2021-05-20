// SPDX-License-Identifier: Apache-2.0

package nftserv

import (
	"encoding/json"
	"fmt"
	"os"
)

const defaultWhitelistedOrigin = "*"

type (
	Config struct {
		AssetsConfig AssetsConfig `json:"assetsConfig"`
		ServerConfig ServerConfig `json:"serverConfig"`
	}

	AssetsConfig struct {
		AssetsPath string `json:"assetsPath"`
		AssetsExt  string `json:"assetsExt"`
	}

	ServerConfig struct {
		Host                 string `json:"host"`
		Port                 uint16 `json:"port"`
		CertFile             string `json:"certFile"`
		KeyFile              string `json:"keyFile"`
		WhitelistedOrigin    string `json:"whitelistedOrigin"`
		MaxTitleLength       int    `json:"maxTitleLength"`
		MaxDescriptionLength int    `json:"maxDescriptionLength"`
	}
)

func ReadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	c := new(Config)
	if err := json.NewDecoder(file).Decode(c); err != nil {
		return nil, fmt.Errorf("decoding server nerd-op config: %w", err)
	}

	if c.ServerConfig.WhitelistedOrigin == "" {
		c.ServerConfig.WhitelistedOrigin = defaultWhitelistedOrigin
	}

	return c, nil
}

// Addr returns the string "{Host}:{Port}"
func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
