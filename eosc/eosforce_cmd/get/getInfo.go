// Copyright © 2018 EOS Canada <info@eoscanada.com>

package get

import (
	"fmt"

	"encoding/json"

	"github.com/fanyang1988/eosc/eosc/eosforce_cmd/utils"
	"github.com/spf13/cobra"
)

var getInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Retrieve blockchain infos, like head block, chain ID, etc..",
	Run: func(cmd *cobra.Command, args []string) {
		api := utils.GetAPI()

		info, err := api.GetInfo()
		utils.ErrorCheck("get info", err)

		data, err := json.MarshalIndent(info, "", "  ")
		utils.ErrorCheck("json marshal", err)

		fmt.Println(string(data))
	},
}

func init() {
	getCmd.AddCommand(getInfoCmd)
}
