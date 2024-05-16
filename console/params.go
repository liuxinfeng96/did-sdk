/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

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
	ParamsFlagTemplateName    = "temp-name"
	ParamsFlagTemplateVersion = "temp-version"
	ParamsFlagTemplatePath    = "temp-path"
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
	ParamsFlagAdminSdkPath    = "admin-sdk-path"
)

var paramsList = map[string]struct {
	shorthand   string
	stringValue string
	usage       string
}{
	ParamsFlagCMSdkPath:       {"C", "", "specify the path of ChainMaker's sdk config file"},
	ParamsFlagAlgorithm:       {"a", "", "specify the public key encryption algorithm. eg. SM2,EC_Secp256k1"},
	ParamsFlagPkPath:          {"p", "", "specify storage path of public key"},
	ParamsFlagPksPath:         {"P", "", "specify the storage path of public key list"},
	ParamsFlagSkPath:          {"s", "", "specify storage path of private key"},
	ParamsFlagSksPath:         {"S", "", "specify the storage path of private key list"},
	ParamsFlagDid:             {"d", "", "specify the did string"},
	ParamsFlagAddress:         {"", "", "specify the address on chain"},
	ParamsFlagDids:            {"D", "", "specify the did list"},
	ParamsFlagListStart:       {"", "", "specify the start index of the query list"},
	ParamsFlagListCount:       {"", "", "specify the number of query list"},
	ParamsFlagListSearch:      {"", "", "specify the keyword of query list"},
	ParamsFlagController:      {"", "", "specify the controller of the DID document"},
	ParamsFlagDocPath:         {"", "", "specify the path of the DID document"},
	ParamsFlagOldDocPath:      {"", "", "specify the path of old DID document"},
	ParamsFlagNewDocPath:      {"", "", "specify the path of updated DID document"},
	ParamsFlagId:              {"", "", "specify ID"},
	ParamsFlagTemplateId:      {"", "", "specify ID of vc template"},
	ParamsFlagTemplateName:    {"", "", "specify name of vc template"},
	ParamsFlagTemplateVersion: {"", "", "specify version of vc template"},
	ParamsFlagTemplatePath:    {"", "", "specify path of vc template"},
	ParamsFlagKeyIndex:        {"", "", "specify the index of the key in the DID document, [1,n]"},
	ParamsFlagSubjectPath:     {"", "", "specify the path of the vc's subject"},
	ParamsFlagExpiration:      {"", "", "specify the validity period of the vc, format [yy-mm-dd]"},
	ParamsFlagVcPath:          {"", "", "specify the path of vc"},
	ParamsFlagType:            {"", "", "specify the type of vc or vp"},
	ParamsFlagIssuer:          {"", "", "specify the issuer's did of vc"},
	ParamsFlagVcList:          {"", "", "specify the path of vc list in vp"},
	ParamsFlagVpPath:          {"", "", "specify the path of vp"},
	ParamsFlagHolder:          {"", "", "specify the holder of vp"},
	ParamsFlagMapKey:          {"", "", "specify the key list of vc template"},
	ParamsFlagMapValue:        {"", "", "specify the value list of vc template"},
	ParamsFlagAdminSdkPath:    {"", "", "specify the path of admin's sdk config file"},
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
