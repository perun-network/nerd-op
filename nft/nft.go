// SPDX-License-Identifier: Apache-2.0

package nft

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
	"github.com/perun-network/erdstall/value"
)

var ErrNotFound = errors.New("NFT not found")

type (
	NFT struct {
		Token common.Address `json:"token"`
		ID    *big.Int       `json:"id"`
		Owner common.Address `json:"owner"`

		// AssetID is the id of the asset in the assets storage. It must be
		// positive. A value of 0 indicates no asset is set (yet).
		AssetID uint `json:"assetId,omitempty"`
		// Secret is set if this is a secret NFT.
		Secret bool   `json:"secret"`
		Title  string `json:"title"`
		Desc   string `json:"desc"`
	}

	Storage interface {
		// Upsert either inserts a new entry into the storage or updates an existing
		// entry. An NFT is identified by the tuple (token, id), so a NFT already
		// exists iff an NFT with the same values for token and id already exists.
		//
		// Field Owner is updated if it is not the zero address.
		// Field AssetID is updated if it is > 0.
		// Field Secret is update if it is true.
		Upsert(nft NFT) error

		// Get gets the NFT identified by token and id from the storage.
		//
		// If it is not found ErrNFTNotFound is returned.
		Get(token common.Address, id *big.Int) (NFT, error)

		// GetAll returns all NFTs in this storage.
		GetAll() ([]NFT, error)
	}
)

// Extract extracts all NFTs from Account `acc` belonging to `owner`.
// The NFT fields Token, ID and Owner will be set whereas ImageID and Secret
// will be set to the default value.
func Extract(owner common.Address, acc tee.Account) (nfts []NFT) {
	for token, val := range acc.Values {
		ids, ok := val.(*value.IDSet)
		if !ok {
			// Skip fungibles
			continue
		}

		for _, id := range *ids {
			nfts = append(nfts, NFT{
				Token: token,
				ID:    id,
				Owner: owner,
			})
		}
	}
	return
}

func (t *NFT) String() string {
	return fmt.Sprintf("NFT{Token: %s, ID: %s, Owner: %s, AssetID: %d, Secret: %t, Title: `%s`, Desc: `%s`}",
		t.Token.String(), t.ID, t.Owner.String(), t.AssetID, t.Secret, t.Title, t.Desc)
}

func (t *NFT) Update(source NFT) {
	if t.Token != source.Token {
		panic("NFT.Update: Token mismatch")
	} else if t.ID.Cmp(source.ID) != 0 {
		panic("NFT.Update: ID mismatch")
	}
	if source.Owner == eth.Zero {
		t.Owner = source.Owner
	}
	if source.AssetID != 0 {
		t.AssetID = source.AssetID
	}
	if source.Secret {
		t.Secret = true
	}
	if source.Title != "" {
		t.Title = source.Title
	}
	if source.Desc != "" {
		t.Desc = source.Desc
	}
}
