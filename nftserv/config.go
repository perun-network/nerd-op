// SPDX-License-Identifier: Apache-2.0

package nftserv

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	AssetsPath string `json:"assetsPath"`
}

func ReadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	c := new(Config)
	return c, json.NewDecoder(file).Decode(c)
}
