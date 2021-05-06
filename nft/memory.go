// SPDX-License-Identifier: Apache-2.0

package nft

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

var _ Storage = (*Memory)(nil)

type Memory struct {
	mu  sync.RWMutex
	mem map[common.Address]map[string]*NFT
}

func NewMemory() *Memory {
	return &Memory{
		mem: make(map[common.Address]map[string]*NFT),
	}
}

func (m *Memory) Upsert(nft NFT) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if exnft, ok := m.get(nft.Token, nft.ID); ok {
		exnft.Update(nft)
		return nil
	}
	m.put(nft)
	return nil
}

func (m *Memory) Get(token common.Address, id *big.Int) (NFT, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nft, ok := m.get(token, id)
	if !ok {
		return NFT{}, ErrNotFound
	}
	return *nft, nil
}

func (m *Memory) get(token common.Address, id *big.Int) (*NFT, bool) {
	tokenNfts, ok := m.mem[token]
	if !ok {
		return nil, false
	}
	nft, ok := tokenNfts[string(id.Bytes())]
	return nft, ok
}

func (m *Memory) put(nft NFT) {
	tokenNfts, ok := m.mem[nft.Token]
	if !ok {
		tokenNfts = make(map[string]*NFT)
		m.mem[nft.Token] = tokenNfts
	}
	tokenNfts[string(nft.ID.Bytes())] = &nft
}

func (m *Memory) TotalSize() (n int) {
	for _, tnfts := range m.mem {
		n += len(tnfts)
	}
	return
}

func (m *Memory) TokenSize(token common.Address) int {
	if tnfts, ok := m.mem[token]; ok {
		return len(tnfts)
	}
	return 0
}
