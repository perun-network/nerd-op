// SPDX-License-Identifier: Apache-2.0

package nftserv

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/nerd-op/asset"
	"github.com/perun-network/nerd-op/nft"
)

type Server struct {
	r      *mux.Router
	nfts   nft.Storage
	assets asset.Storage
}

func New(nftStorage nft.Storage, assetStorage asset.Storage) *Server {
	s := &Server{
		r:      mux.NewRouter(),
		nfts:   nftStorage,
		assets: assetStorage,
	}
	s.r.HandleFunc("/status", s.handleGETstatus).Methods(http.MethodGet)
	const tokenIdSelector = "/{token:0x[0-9a-fA-F]{40}}/{id:[0-9]+}"
	s.r.HandleFunc("/nft"+tokenIdSelector, s.handlePUTnft).Methods(http.MethodPut)
	s.r.HandleFunc("/nft"+tokenIdSelector, s.handleGETnft).Methods(http.MethodGet)
	s.r.HandleFunc("/nft"+tokenIdSelector+"/asset", s.handleGETnftAsset).Methods(http.MethodGet)
	return s
}

// UpdateBalance is the balance handler that can be injected into the operator
// with Operator.OnNewBalance.
func (s *Server) UpdateBalance(owner common.Address, acc tee.Account) {
	nfts := nft.Extract(owner, acc)
	for _, nft := range nfts {
		if err := s.nfts.Upsert(nft); err != nil {
			log.Errorf("Server.UpdateBalance: Error upserting NFT %v: %v", nft, err)
		}
	}
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.r)
}

func (s *Server) ListenAndServeTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, s.r)
}

func (s *Server) handleGETstatus(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}

func (s *Server) handleGETnft(w http.ResponseWriter, r *http.Request) {
	s.handleNFTRequest(w, r, func(tkn nft.NFT) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(tkn); err != nil {
			log.Errorf("Error JSON-marshalling %v: %v", tkn, err)
		}
	})
}

func (s *Server) handleGETnftAsset(w http.ResponseWriter, r *http.Request) {
	s.handleNFTRequest(w, r, func(tkn nft.NFT) {
		ast, err := s.assets.Get(big.NewInt(int64(tkn.AssetID)))
		if err != nil {
			httpError(w, "asset not found: "+err.Error(), http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		if _, err := w.Write(ast); err != nil {
			log.Errorf("Error sending asset for token %v: %v", tkn, err)
		}
	})
}

func (s *Server) handleNFTRequest(w http.ResponseWriter, r *http.Request, handler func(nft.NFT)) {
	var (
		token, id = mustReadTokenID(r)
		tkn, err  = s.nfts.Get(token, id)
	)

	if errors.Is(err, nft.ErrNotFound) {
		httpError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	handler(tkn)
}

func (s *Server) handlePUTnft(w http.ResponseWriter, r *http.Request) {
	var newtkn nft.NFT
	if err := json.NewDecoder(r.Body).Decode(&newtkn); err != nil {
		http.Error(w, "Error decoding token from payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, id := mustReadTokenID(r)
	if token != newtkn.Token || id.Cmp(newtkn.ID) != 0 {
		httpError(w, "Token or ID mismatch between payload and URL", http.StatusBadRequest)
		return
	}

	tkn, err := s.nfts.Get(token, id)
	if err != nil && !errors.Is(err, nft.ErrNotFound) {
		httpError(w, "Error reading existing token: "+err.Error(), http.StatusInternalServerError)
		return
	} else if err == nil && tkn.Owner != eth.Zero && tkn.Owner != newtkn.Owner { // allow fast PUTter for now...
		httpError(w, fmt.Sprintf("Existing token has different owner: %v", tkn), http.StatusConflict)
		return
	}

	s.nfts.Upsert(newtkn)
}

func mustReadTokenID(r *http.Request) (common.Address, *big.Int) {
	var (
		vars            = mux.Vars(r)
		tokenStr, idStr = vars["token"], vars["id"]
		token           = common.HexToAddress(tokenStr)     // valid due to regexp
		id, _           = new(big.Int).SetString(idStr, 10) // can ignore error due to regexp
	)
	return token, id
}

func httpError(w http.ResponseWriter, err string, code int) {
	log.Debugf("Responding error [%d] %v", code, err)
	http.Error(w, err, code)
}
