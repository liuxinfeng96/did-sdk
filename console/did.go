/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-sdk/did"
	"fmt"
	"os"
	"strings"

	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/spf13/cobra"
)

func DidCMD() *cobra.Command {

	didCmd := &cobra.Command{
		Use:   "did",
		Short: "ChainMaker DID did command",
		Long:  "ChainMaker DID did command",
	}

	didCmd.AddCommand(didMethodCMD())
	didCmd.AddCommand(didGenCMD())
	didCmd.AddCommand(didValidCMD())
	didCmd.AddCommand(didGetCMD())
	return didCmd
}

func didMethodCMD() *cobra.Command {

	var sdkPath string

	methodCmd := &cobra.Command{
		Use:   "method",
		Short: "Get did method",
		Long: strings.TrimSpace(
			`Get ChainMaker DID method name from chain.
Example:
$ ./console did method \
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

			method, err := did.GetDidMethodFromChain(c)
			if err != nil {
				return err
			}

			fmt.Printf("the ChainMaker DID method name is: [%s]\n", method)

			return nil
		},
	}

	attachFlagString(methodCmd, ParamsFlagCMSdkPath, &sdkPath)

	return methodCmd
}

func didGenCMD() *cobra.Command {

	var sdkPath, pkPath string

	genDidCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate did string",
		Long: strings.TrimSpace(
			`Generate did string by public key file.
Example:
$ ./console did gen \
--pk-path=./testdata/pk.pem \
--sdk-path=./testdata/sdk_config.yml
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(pkPath) == 0 {
				return ParamsEmptyError(ParamsFlagPkPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			pkPem, err := os.ReadFile(pkPath)
			if err != nil {
				return err
			}

			didStr, err := did.GenerateDidByPK(pkPem, c)
			if err != nil {
				return err
			}

			fmt.Printf("did: [%s]\n", didStr)
			return nil
		},
	}

	attachFlagString(genDidCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(genDidCmd, ParamsFlagPkPath, &pkPath)

	return genDidCmd
}

func didValidCMD() *cobra.Command {

	var sdkPath, didStr string

	isValidCmd := &cobra.Command{
		Use:   "valid",
		Short: "Whether did is valid",
		Long: strings.TrimSpace(
			`Whether did is valid on chain.
Example:
$ ./console did valid \
--did=did:cm:test1 \
--sdk-psath=./testdata/sdk_config.yml 
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(didStr) == 0 {
				return ParamsEmptyError(ParamsFlagDid)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			ok, err := did.IsValidDidOnChain(didStr, c)
			if err != nil || !ok {
				fmt.Printf("whether the did is valid: [%+v]\n, err: [%s]", ok, err.Error())
				return nil
			}

			fmt.Printf("whether the did is valid: [%+v]\n", ok)
			return nil
		},
	}

	attachFlagString(isValidCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(isValidCmd, ParamsFlagDid, &didStr)

	return isValidCmd
}

func didGetCMD() *cobra.Command {
	var sdkPath, pkPath, address string

	getDidCmd := &cobra.Command{
		Use:   "get",
		Short: "Get did string",
		Long: strings.TrimSpace(
			`Get the did string by public key or address. 
Example:
$ ./console did get \
--pk-path=./testdata/pk.pem \
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

			if len(pkPath) == 0 && len(address) == 0 {
				return ParamsEmptyError(ParamsFlagPkPath)
			}

			if len(pkPath) != 0 {
				pkPem, err := os.ReadFile(pkPath)
				if err != nil {
					return err
				}

				didStr, err := did.GetDidByPkFromChain(string(pkPem), c)
				if err != nil {
					return err
				}

				fmt.Printf("did string: [%s]\n", didStr)
				return nil
			}

			if len(address) != 0 {
				didStr, err := did.GetDidByAddressFromChain(address, c)
				if err != nil {
					return err
				}

				fmt.Printf("did string: [%s]\n", didStr)
				return nil
			}

			return nil
		},
	}

	attachFlagString(getDidCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(getDidCmd, ParamsFlagPkPath, &pkPath)
	attachFlagString(getDidCmd, ParamsFlagAddress, &address)

	return getDidCmd
}
