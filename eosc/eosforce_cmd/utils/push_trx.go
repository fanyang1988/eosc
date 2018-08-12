package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eos-go/ecc"
	"github.com/fanyang1988/eos-go/sudo"
	"github.com/fanyang1988/eosc/eosc/fee"
	"github.com/spf13/viper"
)

func permissionToPermissionLevel(in string) (out eos.PermissionLevel, err error) {
	return eos.NewPermissionLevel(in)
}

func permissionsToPermissionLevels(in []string) (out []eos.PermissionLevel, err error) {
	// loop all parameters
	for _, singleArg := range in {

		// if they specified "account@active,account2", handle that too..
		for _, val := range strings.Split(singleArg, ",") {
			level, err := permissionToPermissionLevel(strings.TrimSpace(val))
			if err != nil {
				return out, err
			}

			out = append(out, level)
		}
	}

	return
}

// PushEOSCActions push transaction with actions to eos net
func PushEOSCActions(api *eos.API, actions ...*eos.Action) {
	permissions := viper.GetStringSlice("global-permission")
	if len(permissions) != 0 {
		levels, err := permissionsToPermissionLevels(permissions)
		ErrorCheck("specified --permission(s) invalid", err)

		for _, act := range actions {
			act.Authorization = levels
		}
	}

	opts := &eos.TxOptions{}

	if chainID := viper.GetString("global-offline-chain-id"); chainID != "" {
		opts.ChainID = toSHA256Bytes(chainID, "--offline-chain-id")
	}

	if headBlockID := viper.GetString("global-offline-head-block"); headBlockID != "" {
		opts.HeadBlockID = toSHA256Bytes(headBlockID, "--offline-head-block")
	}

	if err := opts.FillFromChain(api); err != nil {
		fmt.Println("Error fetching tapos + chain_id from the chain:", err)
		os.Exit(1)
	}

	tx := eos.NewTransaction(actions, opts)

	if viper.GetBool("global-sudo-wrap") {
		binTx, err := eos.MarshalBinary(tx)
		ErrorCheck("binary-packing transaction for sudo wrapping", err)

		tx = eos.NewTransaction([]*eos.Action{sudo.NewExec(eos.AccountName("eosio"), eos.HexBytes(binTx))}, opts)
	}

	tx.SetExpiration(time.Duration(viper.GetInt("global-expiration")) * time.Second)

	tx.Fee = fee.GetFeeByActions(actions)

	var signedTx *eos.SignedTransaction
	var packedTx *eos.PackedTransaction

	if !viper.GetBool("global-skip-sign") {
		signKey := viper.GetString("global-offline-sign-key")
		if signKey != "" {
			pubKey, err := ecc.NewPublicKey(signKey)
			ErrorCheck("parsing public key", err)

			api.SetCustomGetRequiredKeys(func(tx *eos.Transaction) ([]ecc.PublicKey, error) {
				return []ecc.PublicKey{pubKey}, nil
			})
		}

		attachWallet(api)

		var err error
		signedTx, packedTx, err = api.SignTransaction(tx, opts.ChainID, eos.CompressionNone)
		ErrorCheck("signing transaction", err)
	} else {
		signedTx = eos.NewSignedTransaction(tx)
	}

	outputTrx := viper.GetString("global-output-transaction")
	if outputTrx != "" {
		printTrx(signedTx, outputTrx)
	} else {
		if packedTx == nil {
			fmt.Println("A signed transaction is required if you want to broadcast it. Remove --skip-sign (or add --output-transaction ?)")
			os.Exit(1)
		}

		isDebug := viper.GetBool("global-debug")
		if isDebug {
			printTrx(signedTx, "")
		}

		// TODO: print the traces
		resp, err := api.PushTransaction(packedTx)
		ErrorCheck("pushing transaction", err)

		//fmt.Println("Transaction submitted to the network. Confirm at https://eosquery.com/tx/" + resp.TransactionID)
		fmt.Println("Transaction submitted to the network. Transaction ID: " + resp.TransactionID)

	}
}

func printTrx(signedTx *eos.SignedTransaction, outputTrx string) {
	cnt, err := json.MarshalIndent(signedTx, "", "  ")
	ErrorCheck("marshalling json", err)

	if outputTrx != "" {
		err = ioutil.WriteFile(outputTrx, cnt, 0644)
	}
	ErrorCheck("writing output transaction", err)
	for _, act := range signedTx.Actions {
		act.SetToServer(false)
	}
	cnt, err = json.MarshalIndent(signedTx, "", "  ")
	ErrorCheck("marshalling json", err)
	fmt.Println(string(cnt))
	fmt.Println("---")
	if outputTrx != "" {
		fmt.Printf("Transaction written to %q\n", outputTrx)
	}
	fmt.Println("Above is a pretty-printed representation of the outputted file")
}
