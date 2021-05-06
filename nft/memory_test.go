// SPDX-License-Identifier: Apache-2.0

package nft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ptest "perun.network/go-perun/pkg/test"

	"github.com/perun-network/erdstall/eth"

	"github.com/perun-network/nerd-op/nft"
	"github.com/perun-network/nerd-op/nft/test"
)

func TestNFTMemory(t *testing.T) {
	var (
		assert = assert.New(t)
		rng    = ptest.Prng(t)
		tkn    = test.NewRandomNFT(rng)
		m      = nft.NewMemory()
	)

	empty, err := m.Get(tkn.Token, tkn.ID)
	assert.Error(err)
	assert.ErrorIs(err, nft.ErrNotFound)
	assert.Equal(empty, nft.NFT{})

	assert.NoError(m.Upsert(tkn))
	get, err := m.Get(tkn.Token, tkn.ID)
	assert.NoError(err)
	assert.Equal(get, tkn)

	tkn.Owner = eth.NewRandomAddress(rng)
	m.Upsert(tkn)
	get, err = m.Get(tkn.Token, tkn.ID)
	assert.NoError(err)
	assert.Equal(get, tkn)
	assert.Equal(m.TotalSize(), 1)
}
