// SPDX-License-Identifier: Apache-2.0

package nftserv

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
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
	s.r.HandleFunc("/nft/{token:0x[0-9a-fA-F]{40}}/{id:[0-9]+}", s.handleGETnft).Methods(http.MethodGet)
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
	var (
		vars            = mux.Vars(r)
		tokenStr, idStr = vars["token"], vars["id"]
		token           = common.HexToAddress(tokenStr)     // valid due to regexp
		id, _           = new(big.Int).SetString(idStr, 10) // can ignore error due to regexp
		tkn, err        = s.nfts.Get(token, id)
	)

	if errors.Is(err, nft.ErrNotFound) {
		httpError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(tkn); err != nil {
		log.Errorf("Error JSON-marshalling %v: %v", tkn, err)
	}
}

func httpError(w http.ResponseWriter, err string, code int) {
	log.Debugf("Responding error [%d] %v", code, err)
	http.Error(w, err, code)
}
