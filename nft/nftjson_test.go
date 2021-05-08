// SPDX-License-Identifier: Apache-2.0

package nft_test

import (
	"testing"

	wiretest "github.com/perun-network/erdstall/wire/test"
	ptest "perun.network/go-perun/pkg/test"

	"github.com/perun-network/nerd-op/nft"
	"github.com/perun-network/nerd-op/nft/test"
)

func TestNFTJSONMarshalling(t *testing.T) {
	rng := ptest.Prng(t)
	tkn := test.NewRandomNFT(rng)
	wiretest.GenericJSONMarshallingTest(t, tkn, new(nft.NFT))
}
