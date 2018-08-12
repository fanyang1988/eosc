// Copyright © 2018 EOS Canada <info@eoscanada.com>

package get

import (
	"github.com/fanyang1988/eosc/eosc/eosforce_cmd"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch information from the blockchain",
}

func init() {
	cmd.RootCmd.AddCommand(getCmd)
}
