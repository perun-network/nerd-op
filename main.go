// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/operator"

	"github.com/perun-network/nerd-op/nft"
	"github.com/perun-network/nerd-op/nftserv"
)

func main() {
	cfgPath := flag.String("config", "config.json", "operator config file path")
	servPath := flag.String("server", "config.json", "NFT server config file path")
	logLevel := flag.String("log-level", "info", "log level")
	flag.Parse()

	lvl, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Main: error parsing log level: %v", err)
	}
	log.SetLevel(lvl)

	cfg, err := operator.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("Main: error reading operator config: %v", err)
	}
	log.Info("Operator config loaded")

	servCfg, err := nftserv.ReadConfig(*servPath)
	if err != nil {
		log.Fatalf("Main: error reading NFT server config: %v", err)
	}
	log.Info("NFT Server config loaded")

	op := operator.SetupWithPrototypeEnclave(cfg, nil)
	go func() {
		if err := op.Serve(cfg.RPCPort); err != nil {
			log.Errorf("Main: Operator.Serve stopped with error %v", err)
		}
	}()

	serv := nftserv.New(nft.NewMemory())
	addr := servCfg.Host + ":" + servCfg.Port
	if err := serv.ListenAndServe(addr); err != nil {
		log.Errorf("Main: NFTServer.ListenAndServe stopped with error %v", err)
	}
}
