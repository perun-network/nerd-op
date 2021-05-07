# NERD: Nifty Erdstall - Operator & NFT Server

This sweet package provides a customized NFT Server alongside the usual Erdstall
operator to server NFT assets. It is part of NERD, started at the ETHGlobal
Scaling 2021 Hackathon.

## Usage

### Building

Build the combined operator and NFT server main with `go build`.

**Warning**: This package currently relies on an internal development version of
Erdstall. It is expected to be found at `../erdstall`. This `replace` directive
is set in `go.mod`.

### Invocation

```
Usage of ./nerd-op:
  -config string
    	operator config file path (default "config.json")
  -log-level string
    	log level (default "info")
  -server string
    	NFT server config file path (default "server.json")
```

### Configuration files

The NERD operator needs two configuration files. One for the operator (default:
`config.json`) and one for the NFT server (default: `server.json`).
Example configurations for a local Ganache setup can be found in folder `demo`.

See [Erdstall](https://github.com/perun-network/erdstall) for further
information on how to configure the operator.

The NFT server serves assets from folder `assetsPath`. The field `assetsExt`
must be set to the extension of all assets in the folder, e.g., `png`. All files
in the assets folder must be of the form `{id}.{ext}`, where id is a base-10
integer. If the files have no extension, `assetsExt` can be omitted or set to
the empty string.

Asset id 0 is reserved to mean that no asset id has been set (yet).

### HTTP API
The NFT server has three endpoints. Field `{token}` is the ERC721 token address.
It must be of the form `0xdead...beef`, i.e., a 20 byte Ethereum hex address.
Field `{id}` is the ID of the NFT on the token contract. It must be a `uint256`
in base 10.

* `GET /nft/{token}/{id}` - returns the current NFT metadata as JSON.
* `PUT /nft/{token}/{id}` - updates the NFT metadata. The payload must contain a
  JSON of the new metadata. See `nft.NFT` for the JSON format. Only the fields
  `assetId` and `secret` can be updated. If the other fields don't match, the
  request errors. Authentication is TBD.
* `GET /nft/{token}/{id}/asset` - returns the NFT's asset as a data stream.

## License
This project is released under the Apache 2.0 license. See LICENSE for further
information.

_Copyright (C) 2021 PolyCrypt GmbH, Darmstadt, Germany_
