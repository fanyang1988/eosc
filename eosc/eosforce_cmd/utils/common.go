package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bronze1man/go-yaml2json"
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eosc/cli"
	eosvault "github.com/fanyang1988/eosc/vault"
	"github.com/spf13/viper"
)

func MustGetWallet() *eosvault.Vault {
	vault, err := setupWallet()
	ErrorCheck("wallet setup", err)
	return vault
}

func setupWallet() (*eosvault.Vault, error) {
	walletFile := viper.GetString("global-vault-file")
	if _, err := os.Stat(walletFile); err != nil {
		return nil, fmt.Errorf("Wallet file %q missing, ", walletFile)
	}

	vault, err := eosvault.NewVaultFromWalletFile(walletFile)
	if err != nil {
		return nil, fmt.Errorf("loading vault, %s", err)
	}

	boxer, err := eosvault.SecretBoxerForType(vault.SecretBoxWrap, viper.GetString("global-kms-gcp-keypath"))
	if err != nil {
		return nil, fmt.Errorf("secret boxer, %s", err)
	}

	if err := vault.Open(boxer); err != nil {
		return nil, err
	}

	return vault, nil
}

func attachWallet(api *eos.API) error {
	walletURLs := viper.GetStringSlice("global-wallet-url")

	if len(walletURLs) != 1 {
		return errors.New("err : Multi-signer not yet implemented or No wallet url")
	}

	// If a `walletURLs` has a Username in the path, use instead of `default`.
	api.SetSigner(eos.NewWalletSigner(eos.New(walletURLs[0]), "default"))
	return nil
}

// GetAPI create eos api
func GetAPI() *eos.API {
	isDebug := viper.GetBool("global-debug")
	res := eos.New(viper.GetString("global-api-url"))
	res.Debug = isDebug
	return res
}

func ErrorCheck(prefix string, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s: %s\n", prefix, err)
		os.Exit(1)
	}
}

func yamlUnmarshal(cnt []byte, v interface{}) error {
	jsonCnt, err := yaml2json.Convert(cnt)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonCnt, v)
}

func loadYAMLOrJSONFile(filename string, v interface{}) error {
	cnt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if strings.HasSuffix(strings.ToLower(filename), ".json") {
		return json.Unmarshal(cnt, v)
	}
	return yamlUnmarshal(cnt, v)
}

func ToAccount(in, field string) eos.AccountName {
	acct, err := cli.ToAccountName(in)
	if err != nil {
		ErrorCheck(fmt.Sprintf("invalid account format for %q", field), err)
	}

	return acct
}

func ToName(in, field string) eos.Name {
	name, err := cli.ToName(in)
	if err != nil {
		ErrorCheck(fmt.Sprintf("invalid name format for %q", field), err)
	}

	return name
}

func toSHA256Bytes(in, field string) eos.SHA256Bytes {
	if len(in) != 64 {
		ErrorCheck(fmt.Sprintf("%q invalid", field), errors.New("should be 64 hexadecimal characters"))
	}

	bytes, err := hex.DecodeString(in)
	ErrorCheck(fmt.Sprintf("invalid hex in %q", field), err)

	return bytes
}
