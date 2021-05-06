// SPDX-License-Identifier: Apache-2.0

package nftserv

type (
	Storage interface {
		Get(id int) ([]byte, error)
	}
)
