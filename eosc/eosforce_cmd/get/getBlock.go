// Copyright © 2018 EOS Canada <info@eoscanada.com>

package get

import (
	"fmt"

	"encoding/json"

	"github.com/fanyang1988/eosc/eosc/eosforce_cmd/utils"
	"github.com/spf13/cobra"
)

var getBlockCmd = &cobra.Command{
	Use:   "block [block id | block height]",
	Short: "Get block data at a given height, or directly with a block hash",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := utils.GetAPI()

		block, err := api.GetBlockByNumOrIDRaw(args[0])
		utils.ErrorCheck("get block", err)

		data, err := json.MarshalIndent(block, "", "  ")
		utils.ErrorCheck("json marshaling", err)

		fmt.Println(string(data))
	},
}

func init() {
	getCmd.AddCommand(getBlockCmd)

	// getBlockCmd.Flags().BoolP("json", "", false, "return producers info in json")

	// for _, flag := range []string{"json"} {
	// 	if err := viper.BindPFlag("get-block-cmd-"+flag, getBlockCmd.Flags().Lookup(flag)); err != nil {
	// 		panic(err)
	// 	}
	// }
}
