/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-sdk/vc"
	"fmt"
	"strings"

	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/spf13/cobra"
)

func VcRevokeCMD() *cobra.Command {

	vcRevokeCmd := &cobra.Command{
		Use:   "vc-revoke",
		Short: "ChainMaker DID vc-revoke command",
		Long:  "ChainMaker DID vc-revoke command",
	}

	vcRevokeCmd.AddCommand(vcRevokeAddCmd())
	vcRevokeCmd.AddCommand(vcRevokeList())

	return vcRevokeCmd
}

func vcRevokeAddCmd() *cobra.Command {

	var idStr, sdkPath string

	vcRevokeAddCmd := &cobra.Command{
		Use:   "add",
		Short: "Add vc revoke list",
		Long: strings.TrimSpace(
			`Add vc revoke list to blockchain.
Example:
$ ./console vc-revoke add \
--id=16516616 \
--sdk-path=./testdata/sdk_config.yml
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {
			if len(idStr) == 0 {
				return ParamsEmptyError(ParamsFlagId)
			}

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			err = vc.RevokeVCOnChain(idStr, c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(vcRevokeAddCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(vcRevokeAddCmd, ParamsFlagId, &idStr)

	return vcRevokeAddCmd
}

func vcRevokeList() *cobra.Command {
	var start, count int
	var search, sdkPath string

	vcRevokeListCmd := &cobra.Command{
		Use:   "list",
		Short: "Get vc revoke list",
		Long: strings.TrimSpace(
			`Get the vc revoke list from blockchain.
Example:
$ ./console vc-revoke list \
--search=111515515 \
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

			list, err := vc.GetVCRevokedListFromChain(search, start, count, c)
			if err != nil {
				return err
			}

			fmt.Printf("get the vc revoke list: [%+v]\n", list)

			return nil
		},
	}

	attachFlagString(vcRevokeListCmd, ParamsFlagListSearch, &search)
	attachFlagString(vcRevokeListCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagInt(vcRevokeListCmd, ParamsFlagListStart, &start)
	attachFlagInt(vcRevokeListCmd, ParamsFlagListCount, &count)

	return vcRevokeListCmd
}
