/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-sdk/did"
	"fmt"
	"strings"

	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/spf13/cobra"
)

func IssuerCMD() *cobra.Command {

	issuerCmd := &cobra.Command{
		Use:   "issuer",
		Short: "ChainMaker DID issuer command",
		Long:  "ChainMaker DID issuer command",
	}

	issuerCmd.AddCommand(issuerAdd())
	issuerCmd.AddCommand(issuerList())
	issuerCmd.AddCommand(issuerDelete())
	return issuerCmd
}

func issuerAdd() *cobra.Command {
	var dids []string
	var sdkPath string

	issuerAddCmd := &cobra.Command{
		Use:   "add",
		Short: "Add issuer did list",
		Long: strings.TrimSpace(
			`Add issuer DID list to blockchain.
Example:
$ ./console issuer add \
--dids=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk_config.yml 
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {
			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(dids) == 0 {
				return ParamsEmptyError(ParamsFlagDids)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			err = did.AddTrustIssuerListToChain(dids, c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(issuerAddCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagStringSlice(issuerAddCmd, ParamsFlagDids, &dids)

	return issuerAddCmd
}

func issuerList() *cobra.Command {
	var start, count int
	var search, sdkPath string

	issuerListCmd := &cobra.Command{
		Use:   "list",
		Short: "Get issuer did list",
		Long: strings.TrimSpace(
			`Get the did list of issuer from blockchain.
Example:
$ ./console issuer list \
--search=did:cm:test1 \
--start=1 \
--count=10 \
--sdk-path=./testdata/sdk_config.yml
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			list, err := did.GetTrustIssuerListFromChain(search, start, count, c)
			if err != nil {
				return err
			}

			fmt.Printf("get the did list of issuer: [%+v]\n", list)

			return nil
		},
	}

	attachFlagString(issuerListCmd, ParamsFlagListSearch, &search)
	attachFlagString(issuerListCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagInt(issuerListCmd, ParamsFlagListStart, &start)
	attachFlagInt(issuerListCmd, ParamsFlagListCount, &count)

	return issuerListCmd
}

func issuerDelete() *cobra.Command {
	var dids []string
	var sdkPath string

	issuerDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete issuer did list",
		Long: strings.TrimSpace(
			`Delete did list of issuer on blockchain.
Example:
$ ./console issuer delete \
--dids=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk_config.yml 
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(dids) == 0 {
				return ParamsEmptyError(ParamsFlagDids)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			err = did.DeleteTrustIssuerListFromChain(dids, c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(issuerDeleteCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagStringSlice(issuerDeleteCmd, ParamsFlagDids, &dids)

	return issuerDeleteCmd
}
