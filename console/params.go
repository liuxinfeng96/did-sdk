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
	ParamsFlagCMSdkPath       = "sdk-path"
	ParamsFlagAlgorithm       = "algo"
	ParamsFlagAlgorithms      = "algos"
	ParamsFlagPkPath          = "pk-path"
	ParamsFlagPksPath         = "pks-path"
	ParamsFlagSkPath          = "sk-path"
	ParamsFlagSksPath         = "sks-path"
	ParamsFlagDid             = "did"
	ParamsFlagAddress         = "address"
	ParamsFlagDids            = "dids"
	ParamsFlagListStart       = "start"
	ParamsFlagListCount       = "count"
	ParamsFlagListSearch      = "search"
	ParamsFlagController      = "controller"
	ParamsFlagDocPath         = "doc-path"
	ParamsFlagOldDocPath      = "old-doc-path"
	ParamsFlagNewDocPath      = "new-doc-path"
	ParamsFlagId              = "id"
	ParamsFlagTemplateId      = "temp-id"
	ParamsFlagTemplateName    = "name"
	ParamsFlagTemplateVersion = "version"
	ParamsFlagTemplatePath    = "template"
	ParamsFlagKeyIndex        = "key-index"
	ParamsFlagSubjectPath     = "subject"
	ParamsFlagExpiration      = "expiration"
	ParamsFlagVcPath          = "vc-path"
	ParamsFlagType            = "type"
	ParamsFlagIssuer          = "issuer"
	ParamsFlagVcList          = "vc-list"
	ParamsFlagVpPath          = "vp-path"
	ParamsFlagHolder          = "holder"
	ParamsFlagMapKey          = "map-key"
	ParamsFlagMapValue        = "map-value"
)

var paramsList = map[string]struct {
	shorthand   string
	stringValue string
	usage       string
}{
	ParamsFlagCMSdkPath:       {"C", "", "specify the path of ChainMaker's sdk config file"},
	ParamsFlagAlgorithm:       {"al", "", "specify the public key encryption algorithm. eg. SM2,EC_Secp256k1"},
	ParamsFlagAlgorithms:      {"als", "", "specify the public key encryption algorithm list. eg. SM2,EC_Secp256k1"},
	ParamsFlagPkPath:          {"pk", "", "specify storage path of public key"},
	ParamsFlagPksPath:         {"pks", "", "specify the storage path of public key list"},
	ParamsFlagSkPath:          {"sk", "", "specify storage path of private key"},
	ParamsFlagSksPath:         {"sks", "", "specify the storage path of private key list"},
	ParamsFlagDid:             {"d", "", "specify the did string"},
	ParamsFlagAddress:         {"a", "", "specify the address on chain"},
	ParamsFlagDids:            {"ds", "", "specify the did list"},
	ParamsFlagListStart:       {"qs", "", "specify the start index of the query list"},
	ParamsFlagListCount:       {"qc", "", "specify the number of query list"},
	ParamsFlagListSearch:      {"qse", "", "specify the keyword of query list"},
	ParamsFlagController:      {"c", "", "specify the controller of the DID document"},
	ParamsFlagDocPath:         {"doc", "", "specify the path of the DID document"},
	ParamsFlagOldDocPath:      {"odoc", "", "specify the path of old DID document"},
	ParamsFlagNewDocPath:      {"ndoc", "", "specify the path of updated DID document"},
	ParamsFlagId:              {"id", "", "specify ID"},
	ParamsFlagTemplateId:      {"ti", "", "specify ID of vc template"},
	ParamsFlagTemplateName:    {"tn", "", "specify name of vc template"},
	ParamsFlagTemplateVersion: {"tv", "", "specify version of vc template"},
	ParamsFlagTemplatePath:    {"temp", "", "specify path of vc template"},
	ParamsFlagKeyIndex:        {"ki", "", "specify the index of the key in the DID document, [1,n]"},
	ParamsFlagSubjectPath:     {"sub", "", "specify the path of the vc's subject"},
	ParamsFlagExpiration:      {"e", "", "specify the validity period of the vc, format [yy-mm-dd]"},
	ParamsFlagVcPath:          {"vc", "", "specify the path of vc"},
	ParamsFlagType:            {"t", "", "specify the type of vc or vp"},
	ParamsFlagIssuer:          {"i", "", "specify the issuer's did of vc"},
	ParamsFlagVcList:          {"vcs", "", "specify the path of vc list in vp"},
	ParamsFlagVpPath:          {"vp", "", "specify the path of vp"},
	ParamsFlagHolder:          {"h", "", "specify the holder of vp"},
	ParamsFlagMapKey:          {"mk", "", "specify the key list of vc template"},
	ParamsFlagMapValue:        {"mv", "", "specify the value list of vc template"},
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

	flags.StringSliceVarP(params, key, f.shorthand, nil, f.usage)
}

func attachFlagInt(cmd *cobra.Command, key string, params *int) {
	flags := cmd.Flags()

	f, ok := paramsList[key]
	if !ok {
		panic("the flag was not found")
	}

	flags.IntVarP(params, key, f.shorthand, 0, f.usage)
}
