package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	ConsoleOutputSuccessfulOperation = "Successful Operation!"
)

func ParamsEmptyError(params string) error {
	return fmt.Errorf("the parameter [%s] cannot be null", params)
}

const (
	ParamsFlagCMSdkPath  = "sdk-path"
	ParamsFlagAlgorithm  = "algo"
	ParamsFlagAlgorithms = "algos"
	ParamsFlagPkPath     = "pk-path"
	ParamsFlagPksPath    = "pks-path"
	ParamsFlagSkPath     = "sk-path"
	ParamsFlagSksPath    = "sks-path"
	ParamsFlagDid        = "did"
	ParamsFlagAddress    = "address"
	ParamsFlagDids       = "dids"
	ParamsFlagListStart  = "start"
	ParamsFlagListCount  = "count"
	ParamsFlagListSearch = "search"
	ParamsFlagController = "controller"
	ParamsFlagDocPath    = "doc-path"
	ParamsFlagOldDocPath = "old-doc-path"
	ParamsFlagNewDocPath = "new-doc-path"
)

var paramsList = map[string]struct {
	shorthand       string
	stringValue     string
	usage           string
	stringValueList []string
	intValue        int
}{
	ParamsFlagCMSdkPath:  {"C", "", "specify the ChainMaker sdk config file path", nil, 0},
	ParamsFlagAlgorithm:  {"al", "", "Specify the public key encryption algorithm. eg. SM2,EC_Secp256k1", nil, 0},
	ParamsFlagAlgorithms: {"als", "", "Specify the public key encryption algorithm list. eg. SM2,EC_Secp256k1", nil, 0},
	ParamsFlagPkPath:     {"pk", "", "specify the public key storage path", nil, 0},
	ParamsFlagPksPath:    {"pks", "", "specify the public key list key storage path", nil, 0},
	ParamsFlagSkPath:     {"sk", "", "specify the private key storage path", nil, 0},
	ParamsFlagSksPath:    {"sks", "", "specify the private key list storage path", nil, 0},
	ParamsFlagDid:        {"d", "", "specify the did string", nil, 0},
	ParamsFlagAddress:    {"a", "", "specify the address on chain", nil, 0},
	ParamsFlagDids:       {"ds", "", "specify the did list", nil, 0},
	ParamsFlagListStart:  {"qs", "", "specify the start index of the query list", nil, 0},
	ParamsFlagListCount:  {"qc", "", "specify the number of query list", nil, 0},
	ParamsFlagListSearch: {"qse", "", "specify the keyword to query the list", nil, 0},
	ParamsFlagController: {"c", "", "specify the controller of the DID document", nil, 0},
	ParamsFlagDocPath:    {"doc", "", "specify the path of the DID document", nil, 0},
	ParamsFlagOldDocPath: {"odoc", "", "specify the path of old DID document", nil, 0},
	ParamsFlagNewDocPath: {"ndoc", "", "specify the path of updated DID document", nil, 0},
}

func attachFlagString(cmd *cobra.Command, key string, params *string) {
	flags := cmd.Flags()

	f, ok := paramsList[key]
	if !ok {
		panic("the flag was not found")
	}

	flags.StringVarP(params, key, f.shorthand, f.stringValue, f.usage)
}

func attachFlagStringSlice(cmd *cobra.Command, key string, params *[]string) {
	flags := cmd.Flags()

	f, ok := paramsList[key]
	if !ok {
		panic("the flag was not found")
	}

	flags.StringSliceVarP(params, key, f.shorthand, f.stringValueList, f.usage)
}

func attachFlagInt(cmd *cobra.Command, key string, params *int) {
	flags := cmd.Flags()

	f, ok := paramsList[key]
	if !ok {
		panic("the flag was not found")
	}

	flags.IntVarP(params, key, f.shorthand, f.intValue, f.usage)
}
