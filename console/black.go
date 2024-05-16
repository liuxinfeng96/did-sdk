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

func BlackCMD() *cobra.Command {

	blackCmd := &cobra.Command{
		Use:   "black",
		Short: "ChainMaker DID black command",
		Long:  "ChainMaker DID black command",
	}

	blackCmd.AddCommand(blackAdd())
	blackCmd.AddCommand(blackList())
	blackCmd.AddCommand(blackDelete())
	return blackCmd
}

func blackAdd() *cobra.Command {
	var dids []string
	var sdkPath string

	blackAddCmd := &cobra.Command{
		Use:   "add",
		Short: "Add did black list",
		Long: strings.TrimSpace(
			`Add did black list DID to blockchain.
Example:
$ ./console black add \
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

			err = did.AddDidBlackListToChain(dids, c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(blackAddCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagStringSlice(blackAddCmd, ParamsFlagDids, &dids)

	return blackAddCmd
}

func blackList() *cobra.Command {
	var start, count int
	var search, sdkPath string

	blackListCmd := &cobra.Command{
		Use:   "list",
		Short: "Get did black list",
		Long: strings.TrimSpace(
			`Get the did black list from blockchain.
Example:
$ ./console black list \
--search=did:cm:test1 \
--start=1 \
--count=10\
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

			list, err := did.GetDidBlackListFromChain(search, start, count, c)
			if err != nil {
				return err
			}

			fmt.Printf("get the did black list: [%+v]\n", list)

			return nil
		},
	}

	attachFlagString(blackListCmd, ParamsFlagListSearch, &search)
	attachFlagString(blackListCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagInt(blackListCmd, ParamsFlagListStart, &start)
	attachFlagInt(blackListCmd, ParamsFlagListCount, &count)

	return blackListCmd
}

func blackDelete() *cobra.Command {
	var dids []string
	var sdkPath string

	blackDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete did black list",
		Long: strings.TrimSpace(
			`Delete did black list DID on blockchain.
Example:
$ ./console black delete \
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

			err = did.DeleteDidBlackListFromChain(dids, c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(blackDeleteCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagStringSlice(blackDeleteCmd, ParamsFlagDids, &dids)

	return blackDeleteCmd
}
