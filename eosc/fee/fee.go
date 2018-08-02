package fee

import (
	"github.com/fanyang1988/eos-go"
)

var actionFee map[eos.ActionName]eos.Asset

func init() {
	actionFee = map[eos.ActionName]eos.Asset{
		eos.ActN("newaccount"): eos.Asset{Amount: 1000, Symbol: eos.EOSSymbol},
		eos.ActN("transfer"):   eos.Asset{Amount: 100, Symbol: eos.EOSSymbol},
	}
}

// GetFeeByAction get fee by action name
func GetFeeByAction(actionName eos.ActionName) eos.Asset {
	fee, ok := actionFee[actionName]
	//fmt.Printf("get key %v %v %v %v", fee.String(), actionName, actionFee, ok)
	if ok {
		//fmt.Printf("get key %v", fee.String())
		return fee
	} else {
		return eos.Asset{Amount: 0, Symbol: eos.EOSSymbol}
	}
}

// GetFeeByActions get fee sum by actions
func GetFeeByActions(actions []*eos.Action) eos.Asset {
	res := eos.Asset{Amount: 0, Symbol: eos.EOSSymbol}
	for _, act := range actions {
		feeAct := GetFeeByAction(act.Name)
		res = res.Add(feeAct)
		//fmt.Printf("add key %v %v\n", res.String(), feeAct.String())
	}
	return res
}
