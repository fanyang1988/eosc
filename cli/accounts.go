package cli

import (
	"fmt"
	"regexp"

	eos "github.com/fanyang1988/eos-go"
)

var reValidAccount = regexp.MustCompile(`[a-z12345]+`)

func ToAccountName(in string) (out eos.AccountName, err error) {
	if !reValidAccount.MatchString(in) {
		err = fmt.Errorf("invalid characters in %q, allowed: 'a' through 'z', and '1', '2', '3', '4', '5'", in)
		return
	}

	if len(in) > 12 {
		err = fmt.Errorf("%q too long, 12 characters allowed maximum", in)
		return
	}

	if len(in) == 0 {
		err = fmt.Errorf("empty")
		return
	}

	return eos.AccountName(in), nil
}

func ToName(in string) (out eos.Name, err error) {
	name, err := ToAccountName(in)
	if err != nil {
		return
	}
	return eos.Name(name), nil
}
