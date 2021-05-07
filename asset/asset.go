// SPDX-License-Identifier: Apache-2.0

package asset

import "math/big"

type (
	Storage interface {
		Get(id *big.Int) ([]byte, error)
	}
)
