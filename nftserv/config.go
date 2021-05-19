// SPDX-License-Identifier: Apache-2.0

package nftserv

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Host       string       `json:"host"`
	Port       uint16       `json:"port"`
	AssetsPath string       `json:"assetsPath"`
	AssetsExt  string       `json:"assetsExt"`
	CertFile   string       `json:"certFile"`
	KeyFile    string       `json:"keyFile"`
	Extra      ServerExtras `json:"extra"`
}

type ServerExtras struct {
	WhitelistedOrigin    *string `json:"whitelistedOrigin"`
	MaxTitleLength       *int    `json:"maxTitleLength"`
	MaxDescriptionLength *int    `json:"maxDescriptionLength"`
}

func ReadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	c := new(Config)
	return c, json.NewDecoder(file).Decode(c)
}

// Addr returns the string "{Host}:{Port}"
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
