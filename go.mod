module github.com/perun-network/nerd-op

go 1.16

require (
	github.com/ethereum/go-ethereum v1.10.1
	github.com/gorilla/mux v1.7.3
	github.com/perun-network/erdstall v0.0.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0
	perun.network/go-perun v0.6.0
)

replace github.com/perun-network/erdstall => ../erdstall
