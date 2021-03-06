// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eos-go/eosforce"
	"github.com/fanyang1988/eosc/eosc/eosforce_cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var transferCmd = &cobra.Command{
	Use:   "transfer [from] [to] [amount]",
	Short: "Transfer from tokens from an account to another",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		from := utils.ToAccount(args[0], "from")
		to := utils.ToAccount(args[1], "to")
		quantity, err := eos.NewEOSAssetFromString(args[2])
		utils.ErrorCheck("invalid amount", err)
		memo := viper.GetString("transfer-cmd-memo")

		api := utils.GetAPI()

		action := eosforce.NewTransfer(from, to, quantity, memo)

		// in eosforce the sys token is use `eosio.transfer` in System to transfer coin
		action.Account = eos.AN("eosio")
		// action.Account = ToAccount(viper.GetString("transfer-cmd-contract"), "--contract")

		_, err = utils.PushEOSCActions(api, action)
		utils.ErrorCheck("push trx err", err)
	},
}

func init() {
	RootCmd.AddCommand(transferCmd)

	transferCmd.Flags().StringP("memo", "m", "", "Memo to attach to the transfer.")
	transferCmd.Flags().StringP("contract", "", "eosio.token", "Contract to send the transfer through. eosio.token is the contract dealing with the native EOS token.")

	for _, flag := range []string{"memo", "contract"} {
		if err := viper.BindPFlag("transfer-cmd-"+flag, transferCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
