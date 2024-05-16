/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-sdk/did"
	"did-sdk/key"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"chainmaker.org/chainmaker/did-contract/model"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/spf13/cobra"
)

func DocCMD() *cobra.Command {

	docCmd := &cobra.Command{
		Use:   "doc",
		Short: "ChainMaker DID doc command",
		Long:  "ChainMaker DID doc command",
	}

	docCmd.AddCommand(docGenCmd())
	docCmd.AddCommand(docAdd())
	docCmd.AddCommand(docGet())
	docCmd.AddCommand(docUpdateLocal())
	docCmd.AddCommand(docUpdate())

	return docCmd
}

func docGenCmd() *cobra.Command {
	var sdkPath, docPath string
	var sksPath, pksPath, controller []string

	docGenCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate did document",
		Long: strings.TrimSpace(
			`Generate the did document.
Example:
$ ./console doc gen \
--sks-path=./testdata/sk.pem \
--pks-path=./testdata/pk.pem \
--controller=did:cm:test1,did:cm:test2 \
--sdk-path=./testdata/sdk_config.yml \
--doc-path=./testdata/doc.json
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(docPath) == 0 {
				return ParamsEmptyError(ParamsFlagDocPath)
			}

			if len(sksPath) == 0 {
				return ParamsEmptyError(ParamsFlagSksPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			if len(sksPath) != len(pksPath) {
				return errors.New("the number of public and private keys names are not equal")
			}

			keyInfos := make([]*key.KeyInfo, 0)

			for i := 0; i < len(sksPath); i++ {
				var pkPem, skPem []byte
				pkPem, err = os.ReadFile(pksPath[i])
				if err != nil {
					return err
				}

				skPem, err = os.ReadFile(sksPath[i])
				if err != nil {
					return err
				}

				keyInfo := &key.KeyInfo{
					PkPEM: pkPem,
					SkPEM: skPem,
				}

				keyInfos = append(keyInfos, keyInfo)

			}

			doc, err := did.GenerateDidDoc(keyInfos, c, controller...)
			if err != nil {
				return err
			}

			err = os.WriteFile(docPath, doc, 0600)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagStringSlice(docGenCmd, ParamsFlagPksPath, &pksPath)
	attachFlagStringSlice(docGenCmd, ParamsFlagSksPath, &sksPath)
	attachFlagString(docGenCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(docGenCmd, ParamsFlagDocPath, &docPath)
	attachFlagStringSlice(docGenCmd, ParamsFlagController, &controller)

	return docGenCmd
}

func docAdd() *cobra.Command {
	var docPath, sdkPath string

	docAddCmd := &cobra.Command{
		Use:   "add",
		Short: "Add did document",
		Long: strings.TrimSpace(
			`Add the did document to blockchain.
Example:
$ ./console doc add \
--doc-path=./testdata/doc.json \
--sdk-path=./testdata/sdk_config.yml 
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {
			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(docPath) == 0 {
				return ParamsEmptyError(ParamsFlagDocPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			doc, err := os.ReadFile(docPath)
			if err != nil {
				return err
			}

			err = did.AddDidDocToChain(string(doc), c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(docAddCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(docAddCmd, ParamsFlagDocPath, &docPath)

	return docAddCmd
}

func docGet() *cobra.Command {
	var didStr, sdkPath, docPath string

	docGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get did document",
		Long: strings.TrimSpace(
			`Get the did document from blockchain.
Example:
$ ./console doc get \
--did=did:cm:test1 \
--sdk-path=./testdata/sdk_config.yml \
--doc-path=./testdata/doc.json
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {
			if len(didStr) == 0 {
				return ParamsEmptyError(ParamsFlagDid)
			}

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(docPath) == 0 {
				return ParamsEmptyError(ParamsFlagDocPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			doc, err := did.GetDidDocFromChain(didStr, c)
			if err != nil {
				return err
			}

			err = os.WriteFile(docPath, doc, 0600)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(docGetCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(docGetCmd, ParamsFlagDocPath, &docPath)
	attachFlagString(docGetCmd, ParamsFlagDid, &didStr)

	return docGetCmd
}

func docUpdateLocal() *cobra.Command {
	var oldDocPath, newDocPath string
	var sksPath, pksPath, controller []string

	docUpdateLocalCmd := &cobra.Command{
		Use:   "update-local",
		Short: "Update did document",
		Long: strings.TrimSpace(
			`Update the did document at local.
Example:
$ ./console doc update-local  \
--sks-path=./testdata/sk.pem \
--pks-path=./testdata/pk.pem \
--controller=did:cm:test1,did:cm:test2 \
--old-doc-path=./testdata/doc.json \
--new-doc-path=./testdata/doc2.json
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {
			if len(oldDocPath) == 0 {
				return ParamsEmptyError(ParamsFlagOldDocPath)
			}

			if len(newDocPath) == 0 {
				return ParamsEmptyError(ParamsFlagNewDocPath)
			}

			if len(sksPath) != len(pksPath) {
				return errors.New("the number of public and private keys names are not equal")
			}

			keyInfos := make([]*key.KeyInfo, 0)

			for i := 0; i < len(sksPath); i++ {
				pkPem, err := os.ReadFile(pksPath[i])
				if err != nil {
					return err
				}

				skPem, err := os.ReadFile(sksPath[i])
				if err != nil {
					return err
				}

				keyInfo := &key.KeyInfo{
					PkPEM: pkPem,
					SkPEM: skPem,
				}

				keyInfos = append(keyInfos, keyInfo)

			}

			oldDocBytes, err := os.ReadFile(oldDocPath)
			if err != nil {
				return err
			}

			var oldDoc model.DidDocument

			err = json.Unmarshal(oldDocBytes, &oldDoc)
			if err != nil {
				return err
			}

			newDoc, err := did.UpdateDidDoc(oldDoc, keyInfos, controller...)
			if err != nil {
				return err
			}

			err = os.WriteFile(newDocPath, newDoc, 0600)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagStringSlice(docUpdateLocalCmd, ParamsFlagPksPath, &pksPath)
	attachFlagStringSlice(docUpdateLocalCmd, ParamsFlagSksPath, &sksPath)
	attachFlagString(docUpdateLocalCmd, ParamsFlagOldDocPath, &oldDocPath)
	attachFlagString(docUpdateLocalCmd, ParamsFlagNewDocPath, &newDocPath)
	attachFlagStringSlice(docUpdateLocalCmd, ParamsFlagController, &controller)

	return docUpdateLocalCmd
}

func docUpdate() *cobra.Command {
	var docPath, sdkPath string

	docUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update did document",
		Long: strings.TrimSpace(
			`Update the did document on blockchain.
Example:
$ ./console doc update \
--doc-path=./testdata/doc.json \
--sdk-path=./testdata/sdk_config.yml 
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {
			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(docPath) == 0 {
				return ParamsEmptyError(ParamsFlagDocPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			doc, err := os.ReadFile(docPath)
			if err != nil {
				return err
			}

			err = did.UpdateDidDocToChain(string(doc), c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(docUpdateCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(docUpdateCmd, ParamsFlagDocPath, &docPath)

	return docUpdateCmd
}
