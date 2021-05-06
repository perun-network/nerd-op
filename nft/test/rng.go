package test

import (
	"math/big"
	"math/rand"

	"github.com/perun-network/erdstall/eth"

	"github.com/perun-network/nerd-op/nft"
)

func NewRandomNFT(rng *rand.Rand) nft.NFT {
	return nft.NFT{
		Token:   eth.NewRandomAddress(rng),
		ID:      big.NewInt(int64(rng.Uint64())),
		Owner:   eth.NewRandomAddress(rng),
		AssetID: uint(rng.Uint32()),
		Secret:  rng.Intn(2) == 1,
	}
}
