/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"did-sdk/vc"
	"errors"
	"fmt"
	"os"
	"strings"

	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/spf13/cobra"
)

func VcTemplateCMD() *cobra.Command {

	vcTemplateCmd := &cobra.Command{
		Use:   "vc-template",
		Short: "ChainMaker DID vc-template command",
		Long:  "ChainMaker DID vc-template command",
	}

	vcTemplateCmd.AddCommand(vcTemplateAddCmd())
	vcTemplateCmd.AddCommand(vcTemplateListCmd())
	vcTemplateCmd.AddCommand(vcTemplateGetCmd())
	vcTemplateCmd.AddCommand(vcTemplateGenCmd())

	return vcTemplateCmd
}

func vcTemplateAddCmd() *cobra.Command {

	var tid, tname, tversion string
	var sdkPath, tempPath string

	vcTemplateAddCmd := &cobra.Command{
		Use:   "add",
		Short: "Add vc template",
		Long: strings.TrimSpace(
			`Add vc template to blockchain.
Example:
$ ./console vc-template add \
--temp-id=1515 \
--temp-name=模板1 \
--temp-version=v1.0.0 \
--temp-path=./testdata/template.json \
--sdk-path=./testdata/sdk_config.yml
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {

			if len(tid) == 0 {
				return ParamsEmptyError(ParamsFlagTemplateId)
			}

			if len(tname) == 0 {
				return ParamsEmptyError(ParamsFlagTemplateName)
			}

			if len(tversion) == 0 {
				return ParamsEmptyError(ParamsFlagTemplateVersion)
			}

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(tempPath) == 0 {
				return ParamsEmptyError(ParamsFlagTemplatePath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			temp, err := os.ReadFile(tempPath)
			if err != nil {
				return err
			}

			err = vc.AddVcTemplateToChain(tid, tname, tversion, temp, c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(vcTemplateAddCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(vcTemplateAddCmd, ParamsFlagTemplateId, &tid)
	attachFlagString(vcTemplateAddCmd, ParamsFlagTemplateName, &tname)
	attachFlagString(vcTemplateAddCmd, ParamsFlagTemplateVersion, &tversion)
	attachFlagString(vcTemplateAddCmd, ParamsFlagTemplatePath, &tempPath)

	return vcTemplateAddCmd
}

func vcTemplateListCmd() *cobra.Command {
	var start, count int
	var search, sdkPath string

	vcTemplateListCmd := &cobra.Command{
		Use:   "list",
		Short: "Get vc template list",
		Long: strings.TrimSpace(
			`Get the vc template list from blockchain.
Example:
$ ./console vc-template list \
--search=模板1 \
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

			list, err := vc.GetVcTemplateListFromChain(search, start, count, c)
			if err != nil {
				return err
			}

			for _, v := range list {
				fmt.Printf("%+v\n", v)
			}

			return nil
		},
	}

	attachFlagString(vcTemplateListCmd, ParamsFlagListSearch, &search)
	attachFlagString(vcTemplateListCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagInt(vcTemplateListCmd, ParamsFlagListStart, &start)
	attachFlagInt(vcTemplateListCmd, ParamsFlagListCount, &count)

	return vcTemplateListCmd
}

func vcTemplateGetCmd() *cobra.Command {
	var tid, sdkPath, tempPath string

	vcTemplateGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get vc template",
		Long: strings.TrimSpace(
			`Get the vc template from blockchain. 
Example:
$ ./console vc-template get \
--temp-id=151515 \
--temp-path=./testdata/template.json \
--sdk-path=./testdata/sdk_config.yml
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {

			if len(tid) == 0 {
				return ParamsEmptyError(ParamsFlagTemplateId)
			}

			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(tempPath) == 0 {
				return ParamsEmptyError(ParamsFlagTemplatePath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			temp, err := vc.GetVcTemplateFromChain(tid, c)
			if err != nil {
				return err
			}

			err = os.WriteFile(tempPath, temp, 0600)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(vcTemplateGetCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(vcTemplateGetCmd, ParamsFlagTemplatePath, &tempPath)
	attachFlagString(vcTemplateGetCmd, ParamsFlagTemplateId, &tid)

	return vcTemplateGetCmd
}

func vcTemplateGenCmd() *cobra.Command {
	var keyList, valueList []string
	var tempPath string

	vcTemplateGenCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate vc template",
		Long: strings.TrimSpace(
			`Generate vc template. 
Example:
$ ./console vc-template gen \
--map-key=name,age,sex \
--map-value=姓名,年龄,性别 \
--temp-path=./testdata/temp.json
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(tempPath) == 0 {
				return ParamsEmptyError(ParamsFlagTemplatePath)
			}

			if len(keyList) == 0 {
				return ParamsEmptyError(ParamsFlagMapKey)
			}

			if len(keyList) != len(valueList) {
				return errors.New("the key list is not equal with the value list")
			}

			keyValueMap := make(map[string]string)

			for i := 0; i < len(keyList); i++ {
				keyValueMap[keyList[i]] = valueList[i]
			}

			tempJson, err := vc.GenerateSimpleVcTemplate(keyValueMap)
			if err != nil {
				return err
			}

			err = os.WriteFile(tempPath, tempJson, 0600)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagStringSlice(vcTemplateGenCmd, ParamsFlagMapKey, &keyList)
	attachFlagStringSlice(vcTemplateGenCmd, ParamsFlagMapValue, &valueList)
	attachFlagString(vcTemplateGenCmd, ParamsFlagTemplatePath, &tempPath)

	return vcTemplateGenCmd
}
