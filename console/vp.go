package main

import (
	"did-sdk/vp"
	"fmt"
	"os"
	"strings"

	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/spf13/cobra"
)

func VpCMD() *cobra.Command {

	vpCmd := &cobra.Command{
		Use:   "vp",
		Short: "ChainMaker DID vp command",
		Long:  "ChainMaker DID vp command",
	}

	vpCmd.AddCommand(vpGenCmd())
	vpCmd.AddCommand(vpVerifyCmd())

	return vpCmd
}

func vpGenCmd() *cobra.Command {
	var algo, skPath, id, holder, vpPath string
	var keyIndex int
	var vpType, vcListPath []string

	vpGenCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate the vp",
		Long: strings.TrimSpace(
			`Generate the vc at local.
Example:
$ ./console vp gen \
--sk-path=./testdata/sk.pem \
--algo=SM2 \
--holder=did:cm:admin \
--id=1231232 \
--vc-list=./testdata/vc1.json,./testdata/vc2.json \
--type=Identity \
--vp-path=./testdata/vp.json
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(skPath) == 0 {
				return ParamsEmptyError(ParamsFlagSkPath)
			}

			if len(algo) == 0 {
				return ParamsEmptyError(ParamsFlagAlgorithm)
			}

			if len(id) == 0 {
				return ParamsEmptyError(ParamsFlagId)
			}

			if len(holder) == 0 {
				return ParamsEmptyError(ParamsFlagHolder)
			}

			if len(vcListPath) == 0 {
				return ParamsEmptyError(ParamsFlagVcList)
			}

			skPem, err := os.ReadFile(skPath)
			if err != nil {
				return err
			}

			vcList := make([]string, 0)

			for _, p := range vcListPath {
				vc, err := os.ReadFile(p)
				if err != nil {
					return err
				}

				vcList = append(vcList, string(vc))
			}

			vp, err := vp.GenerateVP(skPem, algo, keyIndex, holder, id, vcList, vpType...)
			if err != nil {
				return err
			}

			err = os.WriteFile(vpPath, vp, 0777)
			if err != nil {
				return err
			}

			return nil
		},
	}

	attachFlagString(vpGenCmd, ParamsFlagSkPath, &skPath)
	attachFlagString(vpGenCmd, ParamsFlagAlgorithm, &algo)
	attachFlagString(vpGenCmd, ParamsFlagId, &id)
	attachFlagString(vpGenCmd, ParamsFlagHolder, &holder)
	attachFlagString(vpGenCmd, ParamsFlagVpPath, &vpPath)

	attachFlagStringSlice(vpGenCmd, ParamsFlagType, &vpType)
	attachFlagStringSlice(vpGenCmd, ParamsFlagVcList, &vcListPath)

	attachFlagInt(vpGenCmd, ParamsFlagKeyIndex, &keyIndex)

	return vpGenCmd
}

func vpVerifyCmd() *cobra.Command {
	var sdkPath, vpPath string

	vpVerifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify the vp",
		Long: strings.TrimSpace(
			`Verify the vp on blockchain.
Example:
$ ./console vp verify \
--vp-path=./testdata/vp.json \
--sdk-path=./testdata/sdk_config.yml
`,
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(sdkPath) == 0 {
				return ParamsEmptyError(ParamsFlagCMSdkPath)
			}

			if len(vpPath) == 0 {
				return ParamsEmptyError(ParamsFlagVpPath)
			}

			c, err := cmsdk.NewChainClient(cmsdk.WithConfPath(sdkPath))
			if err != nil {
				return err
			}

			vpJson, err := os.ReadFile(vpPath)
			if err != nil {
				return err
			}

			ok, err := vp.VerifyVPOnChain(string(vpJson), c)
			if err != nil {
				return err
			}

			fmt.Printf("the verification result of vp is: [%+v]", ok)

			return nil
		},
	}

	attachFlagString(vpVerifyCmd, ParamsFlagCMSdkPath, &sdkPath)
	attachFlagString(vpVerifyCmd, ParamsFlagVcPath, &vpPath)

	return vpVerifyCmd
}
