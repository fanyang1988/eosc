package main

import (
	// Load all contracts here, so we can always read and decode
	// transactions with those contracts.
	_ "github.com/fanyang1988/eos-go/msig"
	_ "github.com/fanyang1988/eos-go/system"
	_ "github.com/fanyang1988/eos-go/token"

	"github.com/fanyang1988/eosc/eosc/eosforce_cmd"
	_ "github.com/fanyang1988/eosc/eosc/eosforce_cmd/get"
	_ "github.com/fanyang1988/eosc/eosc/eosforce_cmd/vault"
)

var version = "dev"

func init() {
	cmd.Version = version
}

func main() {
	cmd.Execute()
}
