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

	var idStr, tname, tversion string
	var sdkPath, tempPath string

	vcTemplateAddCmd := &cobra.Command{
		Use:   "vc-template",
		Short: "Add vc template",
		Long: strings.TrimSpace(
			`Add vc template to blockchain.
Example:
$ cmc vc-template add \
-id 123333test \
-tn 模板1 \
-tv 1.0.0 \
-temp ./testdata/temp.json \
-C ./testdata/sdk.yaml
`,
		),

		RunE: func(_ *cobra.Command, _ []string) error {

			if len(idStr) == 0 {
				return ParamsEmptyError(ParamsFlagId)
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

			err = vc.AddVcTemplateToChain(idStr, tname, tversion, temp, c)
			if err != nil {
				return err
			}

			fmt.Println(ConsoleOutputSuccessfulOperation)

			return nil
		},
	}

	attachFlagString(vcTemplateAddCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(vcTemplateAddCmd, ParamsFlagId, &idStr)

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
$ cmc vc-template list \
-qse vctemp1 \
-qs 1 \
-qc 10 \
-C ./testdata/sdk.yaml
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

			fmt.Printf("get the vc template list: [%+v]\n", list)

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
	var idStr, sdkPath, tempPath string

	vcTemplateGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get vc template",
		Long: strings.TrimSpace(
			`Get the vc template from blockchain. 
Example:
$ cmc vc-template get \
-id 15152515 \
-C ./testdata/sdk.yaml 
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {

			if len(idStr) == 0 {
				return ParamsEmptyError(ParamsFlagId)
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

			temp, err := vc.GetVcTemplateFromChain(idStr, c)
			if err != nil {
				return err
			}

			err = os.WriteFile(tempPath, temp, 0777)
			if err != nil {
				return err
			}

			return nil
		},
	}

	attachFlagString(vcTemplateGetCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(vcTemplateGetCmd, ParamsFlagTemplatePath, &tempPath)
	attachFlagString(vcTemplateGetCmd, ParamsFlagId, &idStr)

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
$ cmc vc-template gen \ 
-mk name,age,sex \
-mv liu,18,man \
-temp ./testdata/temp.json
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

			err = os.WriteFile(tempPath, tempJson, 0777)
			if err != nil {
				return err
			}

			return nil
		},
	}

	attachFlagStringSlice(vcTemplateGenCmd, ParamsFlagMapKey, &keyList)
	attachFlagStringSlice(vcTemplateGenCmd, ParamsFlagMapValue, &valueList)
	attachFlagString(vcTemplateGenCmd, ParamsFlagTemplatePath, &tempPath)

	return vcTemplateGenCmd
}
