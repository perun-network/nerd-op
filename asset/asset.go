// SPDX-License-Identifier: Apache-2.0

package asset

import "math/big"

type (
	Storage interface {
		Get(id *big.Int) ([]byte, error)
	}

	NoStorage struct{}
)

func (NoStorage) Get(*big.Int) ([]byte, error) { return nil, nil }
