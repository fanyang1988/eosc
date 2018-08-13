package utils

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/cihub/seelog"
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eos-go/ecc"
	"github.com/fanyang1988/eos-go/sudo"
	"github.com/fanyang1988/eosc/eosc/fee"
	"github.com/pkg/errors"
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
func PushEOSCActions(api *eos.API, actions ...*eos.Action) (*eos.PushTransactionFullResp, error) {
	permissions := viper.GetStringSlice("global-permission")
	if len(permissions) != 0 {
		levels, err := permissionsToPermissionLevels(permissions)
		if err != nil {
			return nil, errors.WithMessage(err, "specified --permission(s) invalid")
		}
		for _, act := range actions {
			act.Authorization = levels
		}
	}

	opts := &eos.TxOptions{}

	if err := opts.FillFromChain(api); err != nil {
		return nil, errors.WithMessage(err,
			"Error fetching tapos + chain_id from the chain")
	}

	tx := eos.NewTransaction(actions, opts)

	if viper.GetBool("global-sudo-wrap") {
		binTx, err := eos.MarshalBinary(tx)
		if err != nil {
			return nil, errors.WithMessage(err,
				"binary-packing transaction for sudo wrapping")
		}

		tx = eos.NewTransaction([]*eos.Action{
			sudo.NewExec(
				eos.AccountName("eosio"),
				eos.HexBytes(binTx))}, opts)
	}

	tx.SetExpiration(time.Duration(viper.GetInt("global-expiration")) * time.Second)

	tx.Fee = fee.GetFeeByActions(actions)

	var signedTx *eos.SignedTransaction
	var packedTx *eos.PackedTransaction

	if !viper.GetBool("global-skip-sign") {
		signKey := viper.GetString("global-offline-sign-key")
		if signKey != "" {
			pubKey, err := ecc.NewPublicKey(signKey)
			if err != nil {
				return nil, errors.WithMessage(err, "parsing public key")
			}

			api.SetCustomGetRequiredKeys(func(tx *eos.Transaction) ([]ecc.PublicKey, error) {
				return []ecc.PublicKey{pubKey}, nil
			})
		}

		attachWallet(api)

		var err error
		signedTx, packedTx, err = api.SignTransaction(tx, opts.ChainID, eos.CompressionNone)
		if err != nil {
			return nil, errors.WithMessage(err, "signing trx err")
		}
	} else {
		signedTx = eos.NewSignedTransaction(tx)
	}

	if packedTx == nil {
		return nil, errors.New(
			"A signed transaction is required if you want to broadcast it. " +
				"Remove --skip-sign (or add --output-transaction ?)")
	}

	err := printTrx(signedTx)
	if err != nil {
		seelog.Error("printf Trx err by ", err.Error())
	}

	resp, err := api.PushTransaction(packedTx)
	if err != nil {
		return nil, errors.WithMessage(err, "pushing transaction err")
	}

	return resp, err
}

func printTrx(signedTx *eos.SignedTransaction) error {
	for _, act := range signedTx.Actions {
		act.SetToServer(false)
	}

	cnt, err := json.MarshalIndent(signedTx, "", "  ")
	if err != nil {
		return errors.WithMessage(err, "marshalling json")
	}

	seelog.Trace(string(cnt))
	return nil
}
