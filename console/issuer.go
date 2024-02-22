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
		Use:   "Issuer",
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
$ cmc issuer add -C ./testdata/sdk.yaml -ds did:cm:test1
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
$ cmc issuer list -qse did:cm:test1 -qs 1 -qc 10 -C ./testdata/sdk.yaml
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
$ cmc issuer delete -C ./testdata/sdk.yaml -ds did:cm:test1
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
