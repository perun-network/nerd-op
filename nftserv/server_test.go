// SPDX-License-Identifier: Apache-2.0

package nftserv_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
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

const (
	host = "127.0.0.1"
	port = 13443
	ext  = "pc"
)

func TestServer(t *testing.T) {
	var (
		require             = require.New(t)
		rng                 = ptest.Prng(t)
		nfts                = nft.NewMemory()
		assetsDir           = createTmpAssetsDir(t, ext, 0, 1, 420)
		assets, _           = asset.NewFileStorage(assetsDir)
		defaultServerConfig = nftserv.ServerConfig{
			Host: host,
			Port: port,
		}
		srv        = nftserv.New(nfts, assets, defaultServerConfig)
		owner, acc = randomAccount(rng, 5)
		tv         = acc.Values.OrderedValues()[0]
		ids        = value.MustAsBigInts(tv.Value)
		srverr     = make(chan error, 1)
	)
	assets.SetExtension(ext)

	go func() {
		srverr <- srv.Serve()
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

	// GET /nft/...
	for _, id := range ids {
		resp, err := http.Get(url("nft", tv.Token.String(), id))
		require.NoError(err)
		requireStatus(t, resp, http.StatusOK)
		var tkn nft.NFT
		require.NoError(json.NewDecoder(resp.Body).Decode(&tkn))
		require.Equal(tkn.Token, tv.Token)
		require.Equal(tkn.ID, id)
		require.Equal(tkn.Owner, owner)
		require.Zero(tkn.AssetID)
		require.False(tkn.Secret)
	}

	// GET /nfts
	tkns := nft.Extract(owner, acc)
	resp, err := http.Get(url("nfts"))
	require.NoError(err)
	requireStatus(t, resp, http.StatusOK)
	var getnfts []nft.NFT
	require.NoError(json.NewDecoder(resp.Body).Decode(&getnfts))
	require.Len(getnfts, len(tkns))
	require.ElementsMatch(tkns, getnfts)

	// invalid requests
	expectError := func(geturl string, code int) {
		resp, err := http.Get(geturl)
		require.NoError(err)
		requireStatus(t, resp, code)
	}

	expectError(url("foo"), http.StatusNotFound)
	expectError(url("nft", eth.NewRandomAddress(rng).String(), ids[0]), http.StatusNotFound)
	expectError(url("nft", tv.Token.String(), rng.Int()), http.StatusNotFound)

	// GET /nft/.../asset
	expectAsset := func(token common.Address, id *big.Int, assetId uint) {
		resp, err := http.Get(url("nft", token.String(), id, "asset"))
		require.NoError(err)
		requireStatus(t, resp, http.StatusOK)
		data, err := io.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(strconv.Itoa(int(assetId)), string(data))
	}

	tkn := &tkns[0]
	expectAsset(tkn.Token, tkn.ID, 0)

	// PUT /nft/...
	tkn.AssetID = 420
	resp, err = putAsJSON(url("nft", tkn.Token.String(), tkn.ID), tkn)
	require.NoError(err)
	requireStatus(t, resp, http.StatusOK)
	expectAsset(tkn.Token, tkn.ID, 420)

	tkn.Title = strings.Repeat("pay_respect", 25)
	resp, err = putAsJSON(url("nft", tkn.Token.String(), tkn.ID), tkn)
	require.NoError(err)
	requireStatus(t, resp, http.StatusRequestEntityTooLarge)

	tkn.Title = "Valid Title"
	tkn.Desc = strings.Repeat("fubar", 210)
	resp, err = putAsJSON(url("nft", tkn.Token.String(), tkn.ID), tkn)
	require.NoError(err)
	requireStatus(t, resp, http.StatusRequestEntityTooLarge)
}

func requireStatus(t testing.TB, resp *http.Response, code int) {
	t.Helper()
	if code != resp.StatusCode {
		// only read body if code mismatches; might drain body for later assertion
		// otherwise.
		data, err := io.ReadAll(resp.Body)
		assert.NoErrorf(t, err, "Reading body after status mismatch")
		t.Fatalf("Unexpected status: %s; Body: %s", resp.Status, string(data))
	}
}

func url(elems ...interface{}) string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("http://%s:%v", host, port))
	for _, el := range elems {
		s.WriteString(fmt.Sprintf("/%v", el))
	}
	return s.String()
}

func putAsJSON(url string, obj interface{}) (*http.Response, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return new(http.Client).Do(req)
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

func createTmpAssetsDir(t testing.TB, ext string, assetIds ...uint) string {
	dir := t.TempDir()
	for _, id := range assetIds {
		require.NoError(t, os.WriteFile(
			filepath.Join(dir, fmt.Sprintf("%d.%s", id, ext)),
			[]byte(strconv.Itoa(int(id))),
			0666,
		))
	}
	return dir
}
