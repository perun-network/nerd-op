// SPDX-License-Identifier: Apache-2.0

package nft

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// jsonNFT is an intermediary struct used for custom JSON (Un)Marshalling of
// type NFT.
type jsonNFT struct {
	Token   *common.Address `json:"token"`
	ID      string          `json:"id"`
	Owner   *common.Address `json:"owner"`
	AssetID *uint           `json:"assetId,omitempty"`
	Secret  *bool           `json:"secret"`
}

const idBase = 10

func (t NFT) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonNFT{
		Token:   &t.Token,
		ID:      t.ID.Text(idBase),
		Owner:   &t.Owner,
		AssetID: &t.AssetID,
		Secret:  &t.Secret,
	})
}

func (t *NFT) UnmarshalJSON(data []byte) error {
	jt := jsonNFT{
		Token:   &t.Token,
		Owner:   &t.Owner,
		AssetID: &t.AssetID,
		Secret:  &t.Secret,
	}
	if err := json.Unmarshal(data, &jt); err != nil {
		return fmt.Errorf("unmarshalling into jsonNFT: %w", err)
	}
	t.ID = new(big.Int)
	if _, ok := t.ID.SetString(jt.ID, idBase); !ok {
		return fmt.Errorf("ID value (%s) not a valid base %d number string", jt.ID, idBase)
	}
	return nil
}
