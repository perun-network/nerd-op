// SPDX-License-Identifier: Apache-2.0

package nftserv_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	ptest "perun.network/go-perun/pkg/test"

	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
	"github.com/perun-network/erdstall/value"
	vtest "github.com/perun-network/erdstall/value/test"

	"github.com/perun-network/nerd-op/asset"
	"github.com/perun-network/nerd-op/nft"
	"github.com/perun-network/nerd-op/nftserv"
)

const addr = "127.0.0.1:13443"

func TestServer(t *testing.T) {
	var (
		require    = require.New(t)
		rng        = ptest.Prng(t)
		nfts       = nft.NewMemory()
		srv        = nftserv.New(nfts, asset.NoStorage{})
		srverr     = make(chan error, 1)
		owner, acc = randomAccount(rng, 5)
		tv         = acc.Values.OrderedValues()[0]
		ids        = value.MustAsBigInts(tv.Value)
	)

	go func() {
		srverr <- srv.ListenAndServe(addr)
	}()

	runtime.Gosched()
	// give server time to start listening
	for i := 0; i < 50; i++ {
		t.Logf("Status try %d", i)
		resp, err := http.Get(url("status"))
		if err != nil {
			select {
			case serr := <-srverr:
				t.Fatalf("Server unexpectedly returned with error: %v", serr)
			case <-time.After(2 * time.Millisecond):
				continue
			}
		}
		status, rerr := io.ReadAll(resp.Body)
		require.NoError(rerr)
		resp.Body.Close()
		require.Equal(resp.StatusCode, http.StatusOK)
		require.Equal(string(status), "OK")
		break
	}

	srv.UpdateBalance(owner, acc)

	for _, id := range ids {
		resp, err := http.Get(url("nft", tv.Token.String(), id))
		require.NoError(err)
		require.Equalf(resp.StatusCode, http.StatusOK, "Unexpected status: %s", resp.Status)
		var tkn nft.NFT
		require.NoError(json.NewDecoder(resp.Body).Decode(&tkn))
		require.Equal(tkn.Token, tv.Token)
		require.Equal(tkn.ID, id)
		require.Equal(tkn.Owner, owner)
		require.Zero(tkn.AssetID)
		require.False(tkn.Secret)
	}

	expectError := func(geturl string, code int) {
		resp, err := http.Get(geturl)
		require.NoError(err)
		require.Equalf(resp.StatusCode, code, "Unexpected status: %s", resp.Status)
	}

	expectError(url("foo"), http.StatusNotFound)
	expectError(url("nft", eth.NewRandomAddress(rng).String(), ids[0]), http.StatusNotFound)
	expectError(url("nft", tv.Token.String(), rng.Int()), http.StatusNotFound)
}

func url(elems ...interface{}) string {
	var s strings.Builder
	s.WriteString("http://" + addr)
	for _, el := range elems {
		s.WriteString(fmt.Sprintf("/%v", el))
	}
	return s.String()
}

func randomAccount(rng *rand.Rand, numNFTs int) (common.Address, tee.Account) {
	return eth.NewRandomAddress(rng), tee.Account{
		Nonce: rng.Uint64(),
		Values: value.TokenValues(
			eth.NewRandomAddress(rng),
			vtest.NewRandomIDSet(rng, numNFTs),
		),
	}
}
